package todoist

import (
	"context"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

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

// AddProject defines the options for creating a new project.
type AddProject struct {
	Name       string `json:"name"`
	Color      int64  `json:"color,omitempty"`
	ParentID   int64  `json:"parent_id,omitempty"`
	ChildOrder int64  `json:"child_order,omitempty"`
	IsFavorite int64  `json:"is_favorite,omitempty"`

	TempID string `json:"-"`
}

// Add a new project.
func (s *ProjectsService) Add(ctx context.Context, syncToken string, addProject *AddProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Add")

	id := uuid.New().String()
	tempID := addProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	addCommand := &Command{
		Type:   "project_add",
		Args:   addProject,
		UUID:   id,
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

// UpdateProject defines the options for updating an existing project.
type UpdateProject struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	Color      int64  `json:"color,omitempty"`
	Collapsed  int64  `json:"collapsed,omitempty"`
	IsFavorite int64  `json:"is_favorite,omitempty"`

	TempID string `json:"-"`
}

// Update an existing project.
func (s *ProjectsService) Update(ctx context.Context, syncToken string, updateProject *UpdateProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Update")

	id := uuid.New().String()
	tempID := updateProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	updateCommand := &Command{
		Type:   "project_update",
		Args:   updateProject,
		UUID:   id,
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

type MoveProject struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id"`

	TempID string `json:"-"`
}

// Update parent project relationships of the project.
func (s *ProjectsService) Move(ctx context.Context, syncToken string, moveProject *MoveProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Move")

	id := uuid.New().String()
	tempID := moveProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	moveCommand := &Command{
		Type:   "project_move",
		Args:   moveProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []*Command{moveCommand}

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

type DeleteProject struct {
	ID string `json:"id"`

	TempID string `json:"-"`
}

// Delete an existing project and all its descendants.
func (s *ProjectsService) Delete(ctx context.Context, syncToken string, deleteProject *DeleteProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Delete")

	id := uuid.New().String()
	tempID := deleteProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	deleteCommand := &Command{
		Type:   "project_delete",
		Args:   deleteProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []*Command{deleteCommand}

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

type ArchiveProject struct {
	ID string `json:"id"`

	TempID string `json:"-"`
}

// Archive a project and its descendants.
func (s *ProjectsService) Archive(ctx context.Context, syncToken string, archiveProject *ArchiveProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Archive")

	id := uuid.New().String()
	tempID := archiveProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	archiveCommand := &Command{
		Type:   "project_archive",
		Args:   archiveProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []*Command{archiveCommand}

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

type UnarchiveProject struct {
	ID string `json:"id"`

	TempID string `json:"-"`
}

// Unarchive a project. No ancestors will be unarchived along with
// the unarchived project. Instead, the project is unarchived alone,
// loses any parent relationship (becomes a root project), and is
// placed at the end of the list of other root projects.
func (s *ProjectsService) Unarchive(ctx context.Context, syncToken string, unarchiveProject *UnarchiveProject) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Unarchive")

	id := uuid.New().String()
	tempID := unarchiveProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	unarchiveCommand := &Command{
		Type:   "project_unarchive",
		Args:   unarchiveProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []*Command{unarchiveCommand}

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

type ReorderedProject struct {
	ID         string `json:"id"`
	ChildOrder int32  `json:"child_order"`
}

type ReorderProjects struct {
	Projects []ReorderedProject `json:"projects"`

	TempID string `json:"-"`
}

// The command updates `child_order` properties of items in bulk.
func (s *ProjectsService) Reorder(ctx context.Context, syncToken string, reorderProjects *ReorderProjects) ([]*Project, *CommandResponse, error) {
	s.client.Logln("---------- Projects.Reorder")

	id := uuid.New().String()
	tempID := reorderProjects.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	reorderProjectsCommand := &Command{
		Type:   "project_reorder",
		Args:   reorderProjects,
		UUID:   id,
		TempID: tempID,
	}

	commands := []*Command{reorderProjectsCommand}

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

type ProjectInfo struct {
	Project *Project
	Notes   []interface{} // TODO use the actual notes struct
}

// This function is used to extract detailed information about the project,
// including all the notes. It's especially important because on initial load
// we return no more than the last 10 notes. If a client requires more, they
// can be downloaded using this endpoint. It returns a JSON object with the
// project, and optionally the notes attributes.
func (s *ProjectsService) GetProjectInfo(ctx context.Context, syncToken string, ID string, allData bool) (*ProjectInfo, error) {
	s.client.Logln("---------- Projects.GetProjectInfo")

	s.client.SetDebug(false)
	req, err := s.client.NewRequest(syncToken, []string{}, nil)
	if err != nil {
		return nil, err
	}
	s.client.SetDebug(true)

	// Update the URL
	req.URL, _ = url.Parse(defaultBaseURL + "/projects/get")

	// Parse the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	form, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	// Remove the "commands" form field since we don't use it in this request
	form.Del("commands")

	// Add GetProjectInfo-specific fields
	form.Add("project_id", ID)
	form.Add("all_data", strconv.FormatBool(allData))

	for k := range form {
		s.client.Logf("%-15s %-30s\n", k, form.Get(k))
	}
	s.client.Logln()

	bodyReader := strings.NewReader(form.Encode())

	// Set the updated content-length header or else http/2 will complain about
	// request body being larger than the content length
	req.ContentLength = int64(bodyReader.Len())

	// Add encoded form back to the original request body
	req.Body = io.NopCloser(bodyReader)

	var projectInfoResponse *ProjectInfo
	_, err = s.client.Do(ctx, req, &projectInfoResponse)
	if err != nil {
		return nil, err
	}

	return projectInfoResponse, nil
}
