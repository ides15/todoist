# todoist

## Creating a Client

```go
client := todoist.Client{
    Token: "12345",
    BaseURL: "https://api.todoist.com/sync/v8/sync",
}
```

## Getting All \<Resource>

```go
// Get all projects
projects, err := client.Projects()

// Get all items
items, err := client.Items()
```

## Projects

### Getting a Project

```go
// Get project with the ID of 12345
project, err := client.Project(12345)
```

### Create a Project

```go
project, err := client.AddProject(todoist.Project{
    Name: "Project Name",
    Color: 30,
    ParentID: 0,
    ChildOrder: nil,
    IsFavorite: 1,
})
```

### Updating a Project

```go
projects, err := client.Projects()
project := projects[0]

project.Name = "Updated Project Name")
project.Color = 40

project, err := client.UpdateProject(project)
```

### Moving a Project

```go
projects, err := client.Projects()
project := projects[0]

// Move project under the parent project with the ID of "1234"
err := client.MoveProject(project.ID, 1234)
```

### Deleting a Project

```go
projects, err := client.Projects()
project := projects[0]

err := client.DeleteProject(project.ID)
```

### Archive a Project

```go
projects, err := client.Projects()
project := projects[0]

err := client.ArchiveProject(project.ID)
```

### Unarchive a Project

```go
projects, err := client.Projects()
project := projects[0]

err := client.UnarchiveProject(project.ID)
```

### Reorder a Project

```go
projects, err := client.Projects()
project1 := projects[0]
project2 := projects[1]

err := client.ReorderProjects([]todoist.ProjectReorder{
    {
        ID: project1.ID,
        ChildOrder: 2,
    },
    {
        ID: project2.ID,
        ChildOrder: 1,
    },
})
```
