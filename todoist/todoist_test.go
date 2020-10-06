package todoist

import (
	"bytes"
	"context"
	"log"
	"net/http"
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
	client, err := NewClient(apiToken, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	// List all projects
	projects, _, err := client.Projects.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(*project.ID, *project.Name)
	}

	// Add a new project
	// Specify a TempID if you want to use it in the future, otherwise it will create one for you
	tempID := "e061fa23-524b-4665-9034-05928dc47617"
	projects, resp, err := client.Projects.Add(context.Background(), "", &AddProject{
		Name:   "first new project...",
		TempID: tempID,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(*project.ID, *project.Name)
	}

	// Update the project we just added
	projects, _, err = client.Projects.Update(context.Background(), "", &UpdateProject{
		// get the temp ID of the project we just added so we can update the title
		ID:   strconv.Itoa(int(resp.TempIDMapping[tempID])),
		Name: "an *updated* project!!!",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(*project.ID, *project.Name)
	}
}

func Test_Client_Logging(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	c, _ := NewClient("12345", nil, false)

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

	c1, _ := NewClient("12345", nil, true)

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
	_, err := NewClient("", nil, false)
	if err == nil {
		t.Error("expected an error for an empty API token, got nil")
	}

	c, err := NewClient("12345", nil, false)
	emptyClient := &http.Client{}
	if !reflect.DeepEqual(emptyClient, c.client) {
		t.Errorf("expected http client to be 'emptyClient', got %+v", c.client)
	}

	timeoutHTTPClient := &http.Client{Timeout: 5 * time.Second}
	c1, err := NewClient("12345", timeoutHTTPClient, false)
	if c1.client != timeoutHTTPClient {
		t.Errorf("expected http client to be 'timeoutHTTPClient', got %+v", c1.client)
	}

	if c1.Debug != false {
		t.Errorf("expected client debug flag to be false, got %t", c1.Debug)
	}

	c2, err := NewClient("12345", nil, true)
	if err != nil {
		t.Errorf("expected no error, received %v", err)
	}
	if c2.Debug != true {
		t.Errorf("expected client debug flag to be true, got %t", c2.Debug)
	}
}
