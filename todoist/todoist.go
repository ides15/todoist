package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL = "https://api.todoist.com/sync/v8/sync"
	userAgent      = "todoist-go/1.0.0"
)

// A Client manages communication with the Todoist API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	Debug bool // Flag denoting if debug logging statements should be shown or not

	BaseURL *url.URL // Base URL for API endpoints. Defaults to the public Todoist API (Sync API).

	APIToken string // API Token for authenticating API calls. Found in the Integrations tab of the Todoist user settings.

	UserAgent string // User agent used when communicating with the Todoist API.

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Todoist API.
	Projects *ProjectsService
}

// Logf logs a format string and values to output if the client's debug mode is set to true.
func (c *Client) Logf(format string, a ...interface{}) {
	if c.Debug {
		log.Printf(format, a...)
	}
}

// Logln logs values to output if the client's debug mode is set to true.
func (c *Client) Logln(a ...interface{}) {
	if c.Debug {
		log.Println(a...)
	}
}

type service struct {
	client *Client
}

// NewClient returns a new Todoist API client.
func NewClient(apiToken string, client *http.Client, debug bool) (*Client, error) {
	if apiToken == "" {
		return nil, errors.New("apiToken cannot be empty")
	}

	if client == nil {
		client = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    client,
		Debug:     debug,
		BaseURL:   baseURL,
		APIToken:  apiToken,
		UserAgent: userAgent,
	}

	c.common.client = c

	c.Projects = (*ProjectsService)(&c.common)

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
func (c *Client) NewRequest(syncToken string, resourceTypes []string, commands []*Command) (*http.Request, error) {
	form := url.Values{}

	var resourceTypesStr string
	if len(resourceTypes) != 0 {
		resourceTypes, err := json.Marshal(resourceTypes)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("resourceTypes unable to be serialized: %v", resourceTypes))
		}
		resourceTypesStr = string(resourceTypes)
	} else {
		resourceTypes, _ := json.Marshal([]string{"all"})
		resourceTypesStr = string(resourceTypes)
	}
	form.Add("resource_types", resourceTypesStr)

	if commands != nil && len(commands) != 0 {
		commandsBytes, err := json.Marshal(commands)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("commands unable to be serialized: %v", commands))
		}
		form.Add("commands", string(commandsBytes))
	}

	form.Add("token", c.APIToken)

	if syncToken == "" {
		syncToken = "*"
	}
	form.Add("sync_token", syncToken)

	c.Logln("token\t\t", form.Get("token"))
	c.Logln("sync_token\t\t", form.Get("sync_token"))
	c.Logln("resource_types\t", form.Get("resource_types"))
	c.Logln("commands\t\t", form.Get("commands"))
	c.Logln()

	req, err := http.NewRequest(http.MethodPost, c.BaseURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Response is a Todoist API response. This wraps the standard http.Response
// returned from Todoist.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// TODO: find out if I really need a ReadResponse and CommandResponse, and if I can just combine them.

// ReadResponse is a Todoist API response for a read request.
type ReadResponse struct {
	SyncToken     *string          `json:"sync_token,omitempty"`
	FullSync      *bool            `json:"full_sync,omitempty"`
	TempIDMapping map[string]int64 `json:"temp_id_mapping,omitempty"`

	// user	A user object.
	Projects []*Project `json:"projects,omitempty"`
	// items	An array of item objects.
	// notes	An array of item note objects.
	// project_notes	An array of project note objects.
	// sections	An array of section objects.
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
	TempIDMapping map[string]int64       `json:"temp_id_mapping,omitempty"`
	SyncStatus    map[string]interface{} `json:"sync_status,omitempty"`

	Projects []*Project `json:"projects,omitempty"`
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
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

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// BaseError reports an error caused by a Todoist (sync) API request. BaseError
// will report on any status code response, therefore it is used for non-200 (OK)
// status code responses.
type BaseError struct {
	Response *http.Response `json:"-"` // HTTP response that caused this error

	Tag        *string                `json:"error_tag"`   // error tag
	Code       *int64                 `json:"error_code"`  // error code
	Message    *string                `json:"error"`       // error message
	HTTPCode   *int64                 `json:"http_code"`   // error HTTP code
	ErrorExtra map[string]interface{} `json:"error_extra"` // more detail on errors
}

func (e BaseError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// BadRequestError is used if the request was incorrect.
type BadRequestError struct {
	BaseError
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// UnauthorizedError is used if authentication is required, and has failed, or has not yet been provided.
type UnauthorizedError struct {
	BaseError
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// ForbiddenError is used if the request was valid, but for something that is forbidden.
type ForbiddenError struct {
	BaseError
}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// NotFoundError is used if the requested resource could not be found.
type NotFoundError struct {
	BaseError
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// TooManyRequestsError is used if the user has sent too many requests in a given amount of time.
type TooManyRequestsError struct {
	BaseError
}

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// InternalServerError is used if the request failed due to a server error.
type InternalServerError struct {
	BaseError
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// ServiceUnavailableError is used if the server is currently unable to handle the request.
type ServiceUnavailableError struct {
	BaseError
}

func (e ServiceUnavailableError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", *e.HTTPCode, *e.Tag, *e.Message)
}

// SyncError reports an error caused by a Todoist (sync) API request
// with a 200 (OK) status code response, and contains an embedded BaseError.
// Todoist API docs: https://developer.todoist.com/sync/v8/?shell#response-error
type SyncError struct {
	BaseError // embedded original error

	ID string `json:"-"` // original command UUID
}

// CheckResponse checks the API response for an error, and returns it if
// present. A response is considered an error if it has a status code not equal
// to 200 OK, or it has values in the `sync_status` field that are not equal
// to "ok".
//
// API error responses are expected to have response bodies, and a JSON response
// body that maps to SyncError.
func CheckResponse(r *http.Response) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// TODO: handle this nicer
		return err
	}

	if body == nil {
		return errors.New("response body is nil")
	}

	switch c := r.StatusCode; c {
	// The request was processed successfully.
	// In the Todoist API, a 200 (OK) status code means that the response is at
	// least partially correct. There might be errors in the sync_status field,
	// so we still need to check that field for any errors.
	case http.StatusOK:
		var cr CommandResponse
		// TODO: handle this nicer
		if body != nil {
			err = json.Unmarshal(body, &cr)
			if err != nil {
				// TODO: handle this nicer
				return err
			}
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Range through each of the sync_status values, and map
		// each non "ok" value to a SyncError struct
		for cmdID, cmdResult := range cr.SyncStatus {
			if cmdResult != "ok" {
				// Serialize the command result back into an "unmarshallable" string
				cmdResultBytes, _ := json.Marshal(cmdResult)

				var syncErr SyncError
				err = json.Unmarshal(cmdResultBytes, &syncErr)
				if err != nil {
					return err
				}

				syncErr.ID = cmdID

				return syncErr
			}
		}

		return nil

	// The request was incorrect.
	case http.StatusBadRequest:
		var badRequestError BadRequestError
		if err := json.Unmarshal(body, &badRequestError); err != nil {
			return err
		}

		return badRequestError

	// Authentication is required, and has failed, or has not yet been provided.
	case http.StatusUnauthorized:
		var unauthorizedError UnauthorizedError
		if err := json.Unmarshal(body, &unauthorizedError); err != nil {
			return err
		}

		return unauthorizedError

	// The request was valid, but for something that is forbidden.
	case http.StatusForbidden:
		var forbiddenError ForbiddenError
		if err := json.Unmarshal(body, &forbiddenError); err != nil {
			return err
		}

		return forbiddenError

	// The requested resource could not be found.
	case http.StatusNotFound:
		var notFoundError NotFoundError
		if err := json.Unmarshal(body, &notFoundError); err != nil {
			return err
		}

		return notFoundError

	// The user has sent too many requests in a given amount of time.
	case http.StatusTooManyRequests:
		var tooManyRequestsError TooManyRequestsError
		if err := json.Unmarshal(body, &tooManyRequestsError); err != nil {
			return err
		}

		return tooManyRequestsError

	// The request failed due to a server error.
	case http.StatusInternalServerError:
		var internalServerError InternalServerError
		if err := json.Unmarshal(body, &internalServerError); err != nil {
			return err
		}

		return internalServerError

	// The server is currently unable to handle the request.
	case http.StatusServiceUnavailable:
		var serviceUnavailableError ServiceUnavailableError
		if err := json.Unmarshal(body, &serviceUnavailableError); err != nil {
			return err
		}

		return serviceUnavailableError

	default:
		unknownError := BaseError{
			Response:   r,
			Tag:        pString("UNKNOWN_ERROR"),
			Code:       nil,
			Message:    pString("Unknown error occured."),
			HTTPCode:   pInt64(int64(r.StatusCode)),
			ErrorExtra: map[string]interface{}{},
		}

		return unknownError
	}
}
