# todoist-go

[![Actions Status](https://github.com/ides15/todoist/workflows/Go/badge.svg)](https://github.com/ides15/todoist/actions)

Golang client library for the V8 Todoist Sync API. This repository is in development.

## Installation

```sh
go get -u github.com/ides15/todoist
```

## Creating a Client

All that is required to set up a client is your Todoist API token, which can be found at https://todoist.com/prefs/integrations

```go
client, err := todoist.NewClient("<YOUR_TODOIST_API_TOKEN>", nil)
if err != nil {
    panic(err)
}
```

## Working with Resources

Through `todoist.Client`, you can work with any Todoist resource (projects, notes, items, etc.)

### Getting all projects

```go
projects, _ := client.Projects.GetProjects()

for _, project := range projects {
    fmt.Println(projects)
}
```

### Getting a specific project

```go
project, _ := client.Projects.GetProjectByName("Inbox")

// Or, by ID:
project, _ := client.Projects.GetProjectByID(123)
```

### Creating a project

```go
import "github.com/ides15/todoist/types"

newProject := &types.NewProject{
    Name:       "New Project",
    Color:      1,
    ParentID:   1,
    ChildOrder: 1,
    IsFavorite: 0,
}

err = client.Projects.CreateProject(newProject)
if err != nil {
    t.Fatal(err)
}
```
