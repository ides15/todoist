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

## Working with Resources

Through `todoist.Client`, you can work with any Todoist resource (projects, notes, items, etc).

(See the tests for the most up-to-date examples)

```go
package main

import (
	"context"
	"fmt"

	"github.com/ides15/todoist"
)

func main() {
	client, err := todoist.NewClient("<YOUR_TODOIST_API_TOKEN>")
	if err != nil {
		panic(err)
	}

	projects, _, err := client.Projects.List(context.Background(), "")
	if err != nil {
		panic(err)
	}

	for _, p := range projects {
		fmt.Println(p.ID, p.Name)
	}
}
```
