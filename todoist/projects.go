package todoist

import (
	"context"
	"fmt"
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
	ID             int    `json:"id"`
	LegacyID       *int   `json:"legacy_id"`
	Name           string `json:"name"`
	Color          int    `json:"color"`
	ParentID       *int   `json:"parent_id"`
	LegacyParentID *int   `json:"legacy_parent_id"`
	ChildOrder     int    `json:"child_order"`
	Collapsed      int    `json:"collapsed"`
	Shared         bool   `json:"shared"`
	IsDeleted      int    `json:"is_deleted"`
	IsArchived     int    `json:"is_archived"`
	IsFavorite     int    `json:"is_favorite"`
	SyncID         *int   `json:"sync_id"`
	InboxProject   *bool  `json:"inbox_project"`
	TeamInbox      *bool  `json:"team_inbox"`
}

// List the projects for a user.
func (s *ProjectsService) List(ctx context.Context, syncToken string) ([]Project, ReadResponse, error) {
	s.client.Logln("---------- Projects.List")

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, nil)
	if err != nil {
		return nil, ReadResponse{}, err
	}

	var readResponse ReadResponse
	_, err = s.client.Do(ctx, req, &readResponse)
	if err != nil {
		return nil, readResponse, err
	}

	return readResponse.Projects, readResponse, nil
}

// AddProject defines the options for creating a new project.
type AddProject struct {
	Name       string `json:"name"`
	Color      int    `json:"color,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"`
	ChildOrder int    `json:"child_order,omitempty"`
	IsFavorite int    `json:"is_favorite,omitempty"`

	TempID string `json:"-"`
}

// Add a new project.
func (s *ProjectsService) Add(ctx context.Context, syncToken string, addProject AddProject) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Add")

	id := uuid.New().String()
	tempID := addProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	addCommand := Command{
		Type:   "project_add",
		Args:   addProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{addCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
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
	Color      int    `json:"color,omitempty"`
	Collapsed  int    `json:"collapsed,omitempty"`
	IsFavorite int    `json:"is_favorite,omitempty"`

	TempID string `json:"-"`
}

// Update an existing project.
func (s *ProjectsService) Update(ctx context.Context, syncToken string, updateProject UpdateProject) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Update")

	id := uuid.New().String()
	tempID := updateProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	updateCommand := Command{
		Type:   "project_update",
		Args:   updateProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{updateCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
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
func (s *ProjectsService) Move(ctx context.Context, syncToken string, moveProject MoveProject) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Move")

	id := uuid.New().String()
	tempID := moveProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	moveCommand := Command{
		Type:   "project_move",
		Args:   moveProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{moveCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
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
func (s *ProjectsService) Delete(ctx context.Context, syncToken string, deleteProject DeleteProject) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Delete")

	id := uuid.New().String()
	tempID := deleteProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	deleteCommand := Command{
		Type:   "project_delete",
		Args:   deleteProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{deleteCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
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
func (s *ProjectsService) Archive(ctx context.Context, syncToken string, archiveProject ArchiveProject) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Archive")

	id := uuid.New().String()
	tempID := archiveProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	archiveCommand := Command{
		Type:   "project_archive",
		Args:   archiveProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{archiveCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
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
func (s *ProjectsService) Unarchive(ctx context.Context, syncToken string, unarchiveProject UnarchiveProject) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Unarchive")

	id := uuid.New().String()
	tempID := unarchiveProject.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	unarchiveCommand := Command{
		Type:   "project_unarchive",
		Args:   unarchiveProject,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{unarchiveCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
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
func (s *ProjectsService) Reorder(ctx context.Context, syncToken string, reorderProjects ReorderProjects) ([]Project, CommandResponse, error) {
	s.client.Logln("---------- Projects.Reorder")

	id := uuid.New().String()
	tempID := reorderProjects.TempID
	if tempID == "" {
		tempID = uuid.New().String()
	}

	reorderProjectsCommand := Command{
		Type:   "project_reorder",
		Args:   reorderProjects,
		UUID:   id,
		TempID: tempID,
	}

	commands := []Command{reorderProjectsCommand}

	req, err := s.client.NewRequest(syncToken, []string{"projects"}, commands)
	if err != nil {
		return nil, CommandResponse{}, err
	}

	var commandResponse CommandResponse
	_, err = s.client.Do(ctx, req, &commandResponse)
	if err != nil {
		return nil, commandResponse, err
	}

	return commandResponse.Projects, commandResponse, nil
}

type ProjectInfo struct {
	Project Project       `json:"project"`
	Notes   []interface{} `json:"notes"` // TODO use the actual notes struct
}

// This function is used to extract detailed information about the project,
// including all the notes. It's especially important because on initial load
// we return no more than the last 10 notes. If a client requires more, they
// can be downloaded using this endpoint. It returns a JSON object with the
// project, and optionally the notes attributes.
func (s *ProjectsService) GetProjectInfo(ctx context.Context, syncToken string, ID string, allData bool) (ProjectInfo, error) {
	s.client.Logln("---------- Projects.GetProjectInfo")

	s.client.SetDebug(false)
	req, err := s.client.NewRequest(syncToken, []string{}, nil)
	if err != nil {
		return ProjectInfo{}, err
	}
	s.client.SetDebug(true)

	// Update the URL
	req.URL, _ = url.Parse(defaultBaseURL + "/projects/get")

	// Parse the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ProjectInfo{}, err
	}

	form, err := url.ParseQuery(string(body))
	if err != nil {
		return ProjectInfo{}, err
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

	var projectInfoResponse ProjectInfo
	_, err = s.client.Do(ctx, req, &projectInfoResponse)
	if err != nil {
		return ProjectInfo{}, err
	}

	return projectInfoResponse, nil
}

type ProjectData struct {
	Project  Project       `json:"project"`
	Notes    []interface{} `json:"project_notes"` // TODO use the actual notes struct
	Sections []interface{} `json:"sections"`      // TODO use the actual sections struct
	Items    []interface{} `json:"items"`         // TODO use the actual items struct
}

// Gets a JSON object with the project, its notes, sections and any uncompleted items.
func (s *ProjectsService) GetProjectData(ctx context.Context, syncToken string, projectID string) (ProjectData, error) {
	s.client.Logln("---------- Projects.GetProjectData")

	s.client.SetDebug(false)
	req, err := s.client.NewRequest(syncToken, []string{}, nil)
	if err != nil {
		return ProjectData{}, err
	}
	s.client.SetDebug(true)

	// Update the URL
	req.URL, _ = url.Parse(defaultBaseURL + "/projects/get_data")

	// Parse the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ProjectData{}, err
	}

	form, err := url.ParseQuery(string(body))
	if err != nil {
		return ProjectData{}, err
	}

	// Remove the "commands" form field since we don't use it in this request
	form.Del("commands")

	// Add GetProjectData-specific fields
	form.Add("project_id", projectID)

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

	var projectDataResponse ProjectData
	_, err = s.client.Do(ctx, req, &projectDataResponse)
	if err != nil {
		return ProjectData{}, err
	}

	return projectDataResponse, nil
}

type Pagination struct {
	// The maximum number of archived projects to return (between 1 and 500, default is 500).
	Limit int
	// The offset of the first archived project to return, for pagination purposes (first page is 0).
	Offset int
}

// Get the user's archived projects.
//
// Purposefully leaving `pagination` as a pointer so the caller can optionally pass in
// pagination details. If pagination details are not provided, they are not added to the request.
func (s *ProjectsService) GetArchivedProjects(ctx context.Context, syncToken string, pagination *Pagination) ([]Project, error) {
	s.client.Logln("---------- Projects.GetArchivedProjects")

	s.client.SetDebug(false)
	req, err := s.client.NewRequest(syncToken, []string{}, nil)
	if err != nil {
		return []Project{}, err
	}
	s.client.SetDebug(true)

	// Update the URL
	req.URL, _ = url.Parse(defaultBaseURL + "/projects/get_archived")

	// Parse the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []Project{}, err
	}

	form, err := url.ParseQuery(string(body))
	if err != nil {
		return []Project{}, err
	}

	// Remove the "commands" form field since we don't use it in this request
	form.Del("commands")

	// Add GetProjectData-specific fields
	if pagination != nil {
		form.Add("limit", fmt.Sprint(pagination.Limit))
		form.Add("offset", fmt.Sprint(pagination.Offset))
	}

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

	var archivedProjectsResponse []Project
	_, err = s.client.Do(ctx, req, &archivedProjectsResponse)
	if err != nil {
		return []Project{}, err
	}

	return archivedProjectsResponse, nil
}
