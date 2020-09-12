package todoist_test

import (
	"context"
	"fmt"
	"os"
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
		panic(err)
	}

	// List all projects
	projects, _, err := client.Projects.List(context.Background(), "")
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		fmt.Println(*project.ID, *project.Name)
	}

	// Add a new project
	projects, resp, err := client.Projects.Add(context.Background(), "", &todoist.AddProject{
		Name: "not another new project...",
	})
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		fmt.Println(*project.ID, *project.Name)
	}

	// Update the project we just added
	projects, _, err = client.Projects.Update(context.Background(), "", &todoist.UpdateProject{
		ID:   getFirstFromMap(resp.TempIDMapping), // get the temp ID of the project we just added so we can update the title
		Name: "an *updated* project!!!",
	})
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		fmt.Println(*project.ID, *project.Name)
	}
}
