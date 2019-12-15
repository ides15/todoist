package todoist

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ides15/todoist/types"
)

const (
	defaultBaseURL   = "https://api.todoist.com/sync/v8/sync"
	defaultUserAgent = "todoist-go/1.0.0"
)

type Client struct {
	Token string

	client    *http.Client
	debug     bool
	baseURL   string
	userAgent string

	Projects *ProjectService
}

func NewClient(token string, client *http.Client) (*Client, error) {
	if client == nil {
		client = &http.Client{}
	}

	if token == "" {
		return nil, types.ErrRequiredToken
	}

	c := &Client{
		Token:     token,
		client:    client,
		baseURL:   defaultBaseURL,
		userAgent: defaultUserAgent,
		debug:     false,
	}

	c.Projects = &ProjectService{c: c}

	return c, nil
}

func (c *Client) NewRequest(syncToken string, commands *[]types.Command, resourceTypes *[]string) (*http.Request, error) {
	form := &url.Values{}

	form.Add("token", c.Token)

	if syncToken != "" {
		form.Add("sync_token", syncToken)
	}

	if commands != nil {
		commandsString, _ := json.Marshal(commands)
		form.Add("commands", string(commandsString))
	}

	if resourceTypes != nil {
		resourceTypesString, _ := json.Marshal(resourceTypes)
		form.Add("resource_types", string(resourceTypesString))
	}

	req, err := http.NewRequest("POST", c.baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, types.ErrBuildRequest
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
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

	if res.StatusCode != http.StatusOK {
		errorMessage, err := CreateError(res)
		if err != nil {
			return nil, types.ErrUnknown
		}

		return nil, errorMessage
	}

	return res, nil
}

func (c *Client) Log(v ...interface{}) {
	if c.debug {
		log.Println(v...)
	}
}

func (c *Client) Logf(format string, v ...interface{}) {
	if c.debug {
		log.Printf(format, v...)
	}
}
