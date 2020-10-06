package todoist_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/ides15/todoist/todoist"
)

var (
	// Get the Todoist API token from an environment variable
	apiToken = os.Getenv("TODOIST_API_TOKEN")
)

func getFirstFromMap(m map[string]int64) string {
	for key := range m {
		return key
	}

	return ""
}

func Test_Projects(t *testing.T) {
	// Create the client to interact with Todoist
	client, err := todoist.NewClient(apiToken, true)
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
	tempID := "e061fa23-524b-4665-9034-05928dc47617"
	projects, resp, err := client.Projects.Add(context.Background(), "", &todoist.AddProject{
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
	projects, _, err = client.Projects.Update(context.Background(), "", &todoist.UpdateProject{
		ID:   strconv.Itoa(int(resp.TempIDMapping[tempID])), // get the temp ID of the project we just added so we can update the title
		Name: "an *updated* project!!!",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(*project.ID, *project.Name)
	}
}
