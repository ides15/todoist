# todoist-go

[![Actions Status](https://github.com/ides15/todoist/workflows/Go/badge.svg)](https://github.com/ides15/todoist/actions)

Golang client library for the V8 Todoist Sync API. This repository is in development.

## Installation

```sh
go get github.com/ides15/todoist
```

## Creating a Client

All that is required to set up a client is your Todoist API token, which can be found at https://todoist.com/prefs/integrations

```go
client := todoist.Client{
    Token: "<YOUR_TODOIST_API_TOKEN>",
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
