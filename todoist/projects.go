package todoist

import (
	"context"

	"github.com/google/uuid"
)

// ProjectsService handles communication with the project related
// methods of the Todoist API.
//
// Todoist API docs: https://developer.todoist.com/sync/v8/?shell#projects
type ProjectsService service

// Project represents a Todoist project.
type Project struct {
	ID             *int64  `json:"id,omitempty"`
	LegacyID       *int64  `json:"legacy_id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Color          *int64  `json:"color,omitempty"`
	ParentID       *int64  `json:"parent_id,omitempty"`
	LegacyParentID *int64  `json:"legacy_parent_id,omitempty"`
	ChildOrder     *int64  `json:"child_order,omitempty"`
	Collapsed      *int64  `json:"collapsed,omitempty"`
	Shared         *bool   `json:"shared,omitempty"`
	IsDeleted      *int64  `json:"is_deleted,omitempty"`
	IsArchived     *int64  `json:"is_archived,omitempty"`
	IsFavorite     *int64  `json:"is_favorite,omitempty"`
	SyncID         *int64  `json:"sync_id,omitempty"`
	InboxProject   *bool   `json:"inbox_project,omitempty"`
	TeamInbox      *bool   `json:"team_inbox,omitempty"`
}

// List the projects for a user.
func (s *ProjectsService) List(ctx context.Context, syncToken string) ([]*Project, *ReadResponse, error) {
	s.client.Logln("---------- Projects.List")

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, nil)
	if err != nil {
		return nil, nil, err
	}

	var readResponse *ReadResponse
	_, err = s.client.Do(ctx, req, &readResponse)
	if err != nil {
		return nil, readResponse, err
	}

	return readResponse.Projects, readResponse, nil
}

type AddProject struct {
	Name       string `json:"name"`
	Color      int64  `json:"color"`
	ParentID   int64  `json:"parent_id"`
	ChildOrder int64  `json:"child_order"`
	IsFavorite int64  `json:"is_favorite"`
}

// Add a new project.
func (s *ProjectsService) Add(ctx context.Context, syncToken string, addProject *AddProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Add")

	tempID := uuid.New().String()

	addCommand := &Command{
		Type:   "project_add",
		Args:   addProject,
		UUID:   uuid.New().String(),
		TempID: tempID,
	}

	commands := []*Command{addCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, nil, err
	}

	var commandResponse *CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Projects, commandResponse, nil
}

type UpdateProject struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Color      int64  `json:"color,omitempty"`
	Collapsed  int64  `json:"collapsed,omitempty"`
	IsFavorite int64  `json:"is_favorite,omitempty"`
}

// Update an existing project.
func (s *ProjectsService) Update(ctx context.Context, syncToken string, updateProject *UpdateProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Update")

	tempID := uuid.New().String()

	updateCommand := &Command{
		Type:   "project_update",
		Args:   updateProject,
		UUID:   uuid.New().String(),
		TempID: tempID,
	}

	commands := []*Command{updateCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, nil, err
	}

	var commandResponse *CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Projects, commandResponse, nil
}
