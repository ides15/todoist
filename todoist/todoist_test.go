package todoist

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	// Get the Todoist API token from an environment variable
	apiToken = os.Getenv("TODOIST_API_TOKEN")
)

func Test_Projects(t *testing.T) {
	// Create the client to interact with Todoist
	client, err := NewClient(apiToken)
	if err != nil {
		t.Fatal(err)
	}
	client.SetDebug(false)

	// List all projects
	projects, _, err := client.Projects.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	// Add a new project
	// Specify a TempID if you want to use it in the future, otherwise it will create one for you
	parentProjectTempID := "project1"
	_, resp, err := client.Projects.Add(context.Background(), "", AddProject{
		Name:   "Parent Project",
		TempID: parentProjectTempID,
	})
	if err != nil {
		t.Fatal(err)
	}

	parentProjectID := strconv.Itoa(int(resp.TempIDMapping[parentProjectTempID]))

	childProject1TempID := "project2"
	_, resp, err = client.Projects.Add(context.Background(), "", AddProject{
		Name:   "Child Project 1",
		TempID: childProject1TempID,
	})
	if err != nil {
		t.Fatal(err)
	}

	childProject1ID := strconv.Itoa(int(resp.TempIDMapping[childProject1TempID]))

	childProject2TempID := "project3"
	_, resp, err = client.Projects.Add(context.Background(), "", AddProject{
		Name:   "Child Project 2",
		TempID: childProject2TempID,
	})
	if err != nil {
		t.Fatal(err)
	}

	childProject2ID := strconv.Itoa(int(resp.TempIDMapping[childProject2TempID]))

	// Update the project we just added
	_, _, err = client.Projects.Update(context.Background(), "", UpdateProject{
		// get the temp ID of the project we just added so we can update the title
		ID:   parentProjectID,
		Name: "Updated Project 1",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make project 2 a child of project 1
	_, _, err = client.Projects.Move(context.Background(), "", MoveProject{
		ID:       childProject1ID,
		ParentID: parentProjectID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make project 3 a child of project 1
	_, _, err = client.Projects.Move(context.Background(), "", MoveProject{
		ID:       childProject2ID,
		ParentID: parentProjectID,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.Projects.Reorder(context.Background(), "", ReorderProjects{
		Projects: []ReorderedProject{
			{
				ID:         childProject2ID,
				ChildOrder: 0,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.Projects.Archive(context.Background(), "", ArchiveProject{
		ID: parentProjectID,
	})
	if err != nil {
		t.Fatal(err)
	}

	projects, _, err = client.Projects.Unarchive(context.Background(), "", UnarchiveProject{
		ID: childProject2ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Projects.GetProjectInfo(context.Background(), "", parentProjectID, true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Projects.GetProjectData(context.Background(), "", parentProjectID)
	if err != nil {
		t.Fatal(err)
	}

	archivedProejcts, err := client.Projects.GetArchivedProjects(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}

	archivedProjectsP, err := client.Projects.GetArchivedProjects(context.Background(), "", &Pagination{Limit: 1, Offset: 0})
	if err != nil {
		t.Fatal(err)
	}

	if len(archivedProjectsP) != 1 {
		t.Fatalf("expected the number of archived projects returned to be 1, received %d", len(archivedProjectsP))
	}

	for _, project := range projects {
		if _, _, err = client.Projects.Delete(context.Background(), "", DeleteProject{
			ID: strconv.Itoa(project.ID),
		}); err != nil {
			t.Fatal(err)
		}
	}

	for _, archivedProject := range archivedProejcts {
		if _, _, err = client.Projects.Delete(context.Background(), "", DeleteProject{
			ID: strconv.Itoa(archivedProject.ID),
		}); err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Sections(t *testing.T) {
	client, err := NewClient(apiToken)
	if err != nil {
		t.Fatal(err)
	}
	client.SetDebug(true)

	tempInboxSectionID := "inboxSectionID"
	_, resp, err := client.Sections.Add(context.Background(), "", AddSection{
		Name:         "New Inbox section",
		ProjectID:    2252888543, // Inbox project
		SectionOrder: 0,
		TempID:       tempInboxSectionID,
	})
	if err != nil {
		t.Fatal(err)
	}

	inboxSectionID := strconv.Itoa(int(resp.TempIDMapping[tempInboxSectionID]))

	_, _, err = client.Sections.Update(context.Background(), "", UpdateSection{
		ID:        inboxSectionID,
		Name:      "Updated Inbox section",
		Collapsed: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	sections, _, err := client.Sections.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		t.Logf("%+v\n", section)
	}
}

func Test_Client_Logging(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	c, _ := NewClient("12345")

	// Non-debug Logln
	c.Logln("test")
	cString := buf.String()
	buf.Reset()

	if cString != "" {
		t.Errorf("expected client to not log, got '%s'", cString)
	}

	// Non-debug Logf
	c.Logf("test %s", "case")
	cString = buf.String()
	buf.Reset()

	if cString != "" {
		t.Errorf("expected client to not log, got '%s'", cString)
	}

	c1, _ := NewClient("12345")
	c1.SetDebug(true)

	// Debug Logln
	c1.Logln("test")
	cString = buf.String()
	buf.Reset()

	if !strings.HasSuffix(cString, "test\n") {
		t.Errorf("expected client to log, got '%s'", cString)
	}

	// Debug Logf
	c1.Logf("test %s", "case")
	cString = buf.String()
	buf.Reset()

	if !strings.HasSuffix(cString, "test case\n") {
		t.Errorf("expected client to log, got '%s'", cString)
	}
}

func Test_NewClient(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Error("expected an error for an empty API token, got nil")
	}

	c, _ := NewClient("12345")
	emptyClient := &http.Client{}
	if !reflect.DeepEqual(emptyClient, c.client) {
		t.Errorf("expected http client to be 'emptyClient', got %+v", c.client)
	}

	c1, _ := NewClient("12345")
	timeoutHTTPClient := &http.Client{Timeout: 5 * time.Second}
	c1.SetHTTPClient(timeoutHTTPClient)
	if c1.client != timeoutHTTPClient {
		t.Errorf("expected http client to be 'timeoutHTTPClient', got %+v", c1.client)
	}

	if c1.debug != false {
		t.Error("expected client debug flag to be false, got true")
	}

	c1.SetDebug(true)

	if c1.debug != true {
		t.Error("expected client debug flag to be true, got false")
	}

	c2, err := NewClient("12345")
	c2.SetDebug(true)
	if err != nil {
		t.Errorf("expected no error, received %v", err)
	}
	if c2.debug != true {
		t.Errorf("expected client debug flag to be true, got %t", c2.debug)
	}
}

func Test_NewRequest(t *testing.T) {
	c, err := NewClient("12345")
	if err != nil {
		t.Errorf("expcted no error, received %v", err)
	}

	req, err := c.NewRequest("", []string{}, nil)
	if err != nil {
		t.Errorf("expcted no error, received %v", err)
	}

	syncTokenFormValue := req.FormValue("sync_token")
	if syncTokenFormValue != "*" {
		t.Errorf("sync_token should default to \"*\", received %s", syncTokenFormValue)
	}

	resourceTypesFormValue := req.FormValue("resource_types")
	if resourceTypesFormValue != "[\"all\"]" {
		t.Errorf("resource_types should default to [\"all\"], received %s", resourceTypesFormValue)
	}

	if commandsFormValue, exists := req.Form["commands"]; exists {
		t.Errorf("commands should not be included in form, received %s", commandsFormValue)
	}

	tokenFormValue := req.FormValue("token")
	if tokenFormValue != c.APIToken {
		t.Errorf("token should be %s, received %s", c.APIToken, tokenFormValue)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type header must be application/x-www-form-urlencoded, received %s", contentType)
	}

	userAgent := req.Header.Get("User-Agent")
	if userAgent != c.userAgent {
		t.Errorf("User-Agent should be %s, received %s", c.userAgent, userAgent)
	}

	c.userAgent = ""
	req, _ = c.NewRequest("", []string{"projects"}, []Command{
		{
			Type:   "command_type",
			Args:   "args",
			UUID:   "uuid",
			TempID: "temp_id",
		},
	})

	resourceTypesFormValue = req.FormValue("resource_types")
	if resourceTypesFormValue != "[\"projects\"]" {
		t.Errorf("resource_types JSONified incorrectly, received %s", resourceTypesFormValue)
	}

	_, exists := req.Form["commands"]
	if !exists {
		t.Error("commands expected in form, but were not included")
	} else {
		commandsFormValue := req.Form.Get("commands")
		if commandsFormValue != `[{"type":"command_type","args":"args","uuid":"uuid","temp_id":"temp_id"}]` {
			t.Errorf("commands JSONified incorrectly, received %s", commandsFormValue)
		}
	}

	if userAgent, exists := req.Header["User-Agent"]; exists {
		t.Errorf("User-Agent should not be set in request, received %s", userAgent)
	}

	_, err = c.NewRequest("", []string{"all"}, []Command{
		{
			Type:   "type",
			Args:   c.client, // Just need something that is unserializable
			UUID:   "uuid",
			TempID: "temp_id",
		},
	})
	if err == nil {
		t.Error("expected err serializing commands, received nil")
	}

	c.BaseURL = &url.URL{Host: "localhost#bad-url"}
	_, err = c.NewRequest("", []string{"all"}, nil)
	if err == nil {
		t.Errorf("expected err in new request, received nil")
	}
}
