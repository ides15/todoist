package todoist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL = "https://api.todoist.com/sync/v8"
	userAgent      = "todoist-go/1.0.0"
)

// A Client manages communication with the Todoist API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	debug bool // Flag denoting if debug logging statements should be shown or not

	BaseURL *url.URL // Base URL for API endpoints. Defaults to the public Todoist API (Sync API).

	APIToken string // API Token for authenticating API calls. Found in the Integrations tab of the Todoist user settings.

	userAgent string // User agent used when communicating with the Todoist API.

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Todoist API.
	Projects *ProjectsService
	Sections *SectionsService
	Tasks    *TasksService
}

// Logf logs a format string and values to output if the client's debug mode is set to true.
func (c *Client) Logf(format string, a ...interface{}) {
	if c.debug {
		log.Printf(format, a...)
	}
}

// Logln logs values to output if the client's debug mode is set to true.
func (c *Client) Logln(a ...interface{}) {
	if c.debug {
		log.Println(a...)
	}
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

type service struct {
	client *Client
}

// NewClient returns a new Todoist API client.
func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, errors.New("apiToken cannot be empty")
	}

	baseURL, _ := url.Parse(defaultBaseURL + "/sync")

	c := &Client{
		client:    &http.Client{},
		BaseURL:   baseURL,
		APIToken:  apiToken,
		userAgent: userAgent,
		debug:     false,
	}

	c.common.client = c

	// c.Projects = (*ProjectsService)(&c.common)
	c.Projects = &ProjectsService{client: c}
	c.Sections = &SectionsService{client: c}
	c.Tasks = &TasksService{client: c}

	return c, nil
}

// Command is a Todoist API request parameter for writing Todoist resources.
type Command struct {
	Type   string      `json:"type"`
	Args   interface{} `json:"args"`
	UUID   string      `json:"uuid"`
	TempID string      `json:"temp_id"`
}

// NewRequest creates an API request. If specified, the value pointed to
// by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(syncToken string, resourceTypes []string, commands []Command) (*http.Request, error) {
	form := url.Values{}

	if syncToken == "" {
		syncToken = "*"
	}
	form.Add("sync_token", syncToken)

	if len(resourceTypes) == 0 {
		resourceTypes = []string{"all"}
	}
	resourceTypesBytes, _ := json.Marshal(resourceTypes)
	resourceTypesStr := string(resourceTypesBytes)
	form.Add("resource_types", resourceTypesStr)

	if len(commands) != 0 {
		commandsBytes, err := json.Marshal(commands)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to serialize commands: %v", commands))
		}
		commandsStr := string(commandsBytes)
		form.Add("commands", commandsStr)
	}

	form.Add("token", c.APIToken)

	for k := range form {
		c.Logf("%-15s %-30s\n", k, form.Get(k))
	}
	c.Logln()

	req, err := http.NewRequest(http.MethodPost, c.BaseURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

// TODO: find out if I really need a ReadResponse and CommandResponse, and if I can just combine them.

// ReadResponse is a Todoist API response for a read request.
type ReadResponse struct {
	FullSync      bool           `json:"full_sync"`
	SyncToken     string         `json:"sync_token"`
	TempIDMapping map[string]int `json:"temp_id_mapping"`

	Projects []Project `json:"projects"`
	Sections []Section `json:"sections"`
	Tasks    []Task    `json:"items"`
	// notes	An array of item note objects.
	// project_notes	An array of project note objects.
	// labels	An array of label objects.
	// filters	A array of filter objects.
	// day_orders	A JSON object specifying the order of items in daily agenda.
	// reminders	An array of reminder objects.
	// collaborators	A JSON object containing all collaborators for all shared projects. The projects field contains the list of all shared projects, where the user acts as one of collaborators.
	// collaborators_states	An array specifying the state of each collaborator in each project. The state can be invited, active, inactive, deleted.
	// live_notifications	An array of live_notification objects
	// live_notifications_last_read	What is the last live notification the user has seen? This is used to implement unread notifications.
	// user_settings	A JSON object containing user settings.
	// user_plan_limits	A JSON object containing user plan limits.
}

// CommandResponse is a Todoist API response for a command request.
type CommandResponse struct {
	FullSync      bool                   `json:"full_sync"`
	SyncToken     string                 `json:"sync_token"`
	SyncStatus    map[string]interface{} `json:"sync_status"`
	TempIDMapping map[string]int         `json:"temp_id_mapping"`

	Projects []Project `json:"projects"`
	Sections []Section `json:"sections"`
	Tasks    []Task    `json:"items"`
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer resp.Body.Close()

	err = checkResponseForErrors(resp, v)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			if _, err = io.Copy(w, resp.Body); err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}
