package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

const (
	DefaultBaseURL   = "https://api.todoist.com/sync/v8/sync"
	DefaultUserAgent = "todoist-go/1.0.0"
)

type Client struct {
	*http.Client

	Token     string
	BaseURL   string
	UserAgent string
}

func (c *Client) API() *http.Client {
	if c.Client == nil {
		return http.DefaultClient
	}

	return c.Client
}

func (c *Client) URL() string {
	if c.BaseURL == "" {
		return DefaultBaseURL
	}

	return c.BaseURL
}

func (c *Client) Agent() string {
	if c.UserAgent == "" {
		return DefaultUserAgent
	}

	return c.UserAgent
}

func (c *Client) BuildReq(ctx context.Context, syncToken string, commands *[]Command, resourceTypes *[]string) (*http.Request, error) {
	form := url.Values{}

	form.Add("token", c.Token)

	if syncToken != "" {
		form.Add("sync_token", syncToken)
	}

	if commands != nil {
		commandsString, err := json.Marshal(commands)
		if err != nil {
			return nil, err
		}

		form.Add("commands", string(commandsString))
	}

	if resourceTypes != nil {
		resourceTypesString, err := json.Marshal(resourceTypes)
		if err != nil {
			return nil, err
		}

		form.Add("resource_types", string(resourceTypesString))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.URL(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c *Client) GetProjects(ctx context.Context, syncToken string) (*[]Project, error) {
	req, err := c.BuildReq(ctx, syncToken, nil, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	return res.Projects, nil
}

func (c *Client) GetProjectByID(ctx context.Context, id int, syncToken string) (*Project, error) {
	req, err := c.BuildReq(ctx, syncToken, nil, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	if res.Projects != nil {
		for _, project := range *res.Projects {
			if project.ID == id {
				return &project, nil
			}
		}

		// If the project ID was not found return a not found error
		return nil, errors.New("project not found")
	}

	// If res.Projects is nil for some reason (which shouldn't happen)
	return nil, errors.New("todoist didn't return projects")
}

func (c *Client) GetProjectByName(ctx context.Context, name string, syncToken string) (*Project, error) {
	req, err := c.BuildReq(ctx, syncToken, nil, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	if res.Projects != nil {
		for _, project := range *res.Projects {
			if project.Name == name {
				return &project, nil
			}
		}

		// If the project ID was not found return a not found error
		return nil, errors.New("project not found")
	}

	// If res.Projects is nil for some reason (which shouldn't happen)
	return nil, errors.New("todoist didn't return projects")
}

func (c *Client) AddProject(ctx context.Context, newProject *NewProject, tempID *string, syncToken string) (*Project, error) {
	var temp string
	if tempID == nil {
		temp = uuid.New().String()
	} else {
		temp = *tempID
	}

	if newProject == nil {
		return nil, errors.New("must provide a new project")
	}

	req, err := c.BuildReq(ctx, syncToken, &[]Command{{
		Type:   "project_add",
		Args:   newProject,
		UUID:   uuid.New().String(),
		TempID: temp,
	}}, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	if res.Projects != nil {
		for _, project := range *res.Projects {
			// check if tempIDMapping has this key
			if project.ID == res.TempIDMapping[temp] {
				return &project, nil
			}
		}

		return nil, errors.New("project creation was not successful")
	}

	return nil, errors.New("todoist didn't return projects")
}

func (c *Client) UpdateProject(ctx context.Context, updatedProject *UpdatedProject, syncToken string) (*Project, error) {
	if updatedProject == nil {
		return nil, errors.New("must provide an updated project")
	}

	req, err := c.BuildReq(ctx, syncToken, &[]Command{{
		Type: "project_update",
		Args: updatedProject,
		UUID: uuid.New().String(),
	}}, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	if res.Projects != nil {
		for _, project := range *res.Projects {
			// check if tempIDMapping has this key
			if project.ID == updatedProject.ID {
				return &project, nil
			}
		}

		return nil, errors.New("project update was not successful")
	}

	return nil, errors.New("todoist didn't return projects")
}

func (c *Client) MoveProject(ctx context.Context, moveProject *MovedProject, syncToken string) error {
	req, err := c.BuildReq(ctx, syncToken, &[]Command{{
		Type: "project_move",
		Args: moveProject,
		UUID: uuid.New().String(),
	}}, nil)
	if err != nil {
		return err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (c *Client) DeleteProject(ctx context.Context, deleteProject *DeletedProject, syncToken string) error {
	req, err := c.BuildReq(ctx, syncToken, &[]Command{{
		Type: "project_delete",
		Args: deleteProject,
		UUID: uuid.New().String(),
	}}, nil)
	if err != nil {
		return err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func (c *Client) ReorderProject(ctx context.Context, reorderProject *ReorderedProject, syncToken string) error {
	req, err := c.BuildReq(ctx, syncToken, &[]Command{{
		Type: "project_reorder",
		Args: reorderProject,
		UUID: uuid.New().String(),
	}}, nil)
	if err != nil {
		return err
	}

	response, err := c.API().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
