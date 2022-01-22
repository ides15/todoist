# todoist-go

[![Actions Status](https://github.com/ides15/todoist/workflows/Go/badge.svg)](https://github.com/ides15/todoist/actions)

Golang client library for the V8 Todoist Sync API. This repository is in development.

Influenced by Google's Github API Golang Client: https://github.com/google/go-github

## Installation

```sh
go get -u github.com/ides15/todoist
```

## Creating a Client

All that is required to set up a client is your Todoist API token, which can be found at https://todoist.com/prefs/integrations

```go
client, err := todoist.NewClient("<YOUR_TODOIST_API_TOKEN>")
if err != nil {
    panic(err)
}
```

---

<br/>

## Working with Resources

Through `todoist.Client`, you can work with any Todoist resource (projects, notes, items, etc.)

### Projects

(see the test files for the most up-to-date examples)

```go
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

func Test_Projects(t *testing.T) {
	// Create the client to interact with Todoist
	client, err := todoist.NewClient(apiToken)
	if err != nil {
		t.Fatal(err)
	}

	// List all projects
	projects, _, err := client.Projects.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(project.ID, project.Name)
	}

	// Add a new project
	// Specify a TempID if you want to use it in the future, otherwise it will create one for you
	tempID := "e061fa23-524b-4665-9034-05928dc47617"
	projects, resp, err := client.Projects.Add(context.Background(), "", todoist.AddProject{
		Name:   "first new project...",
		TempID: tempID,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(project.ID, project.Name)
	}

	// Update the project we just added
	projects, _, err = client.Projects.Update(context.Background(), "", todoist.UpdateProject{
		// get the temp ID of the project we just added so we can update the title
		ID:   strconv.Itoa(int(resp.TempIDMapping[tempID])),
		Name: "an *updated* project!!!",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, project := range projects {
		t.Log(project.ID, project.Name)
	}
}
```
