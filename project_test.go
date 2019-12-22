package todoist_test

import (
	"testing"

	"github.com/ides15/todoist/types"
)

func TestGetProjects(t *testing.T) {
	// TODO: check that projects equals what we're expecting
	Setup()

	TestClient.BaseURL = TestServer.URL
	projects, err := TestClient.Projects.GetProjects()
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}

	expectedLength := 2
	if len(projects) != expectedLength {
		t.Fatalf("expected %d projects, received %d", expectedLength, len(projects))
	}
}

func TestGetProjectsBadRequest(t *testing.T) {
	// TODO: check that projects equals what we're expecting
	Setup()

	TestClient.BaseURL = "\t"
	_, err := TestClient.Projects.GetProjects()
	if err == nil {
		t.Fatalf("expected error, received nil")
	}
}

// TODO: add ability to add context into individual resource requests
// func TestGetProjectsDoRequestContextCancel(t *testing.T) {
// 	Setup()
//
// 	TestClient.BaseURL = TestServer.URL
// 	d := time.Now().Add(1 * time.Second)
// 	ctx, cancel := context.WithDeadline(context.Background(), d)
// 	cancel()
//
// 	_, err := TestClient.Projects.GetProjects(ctx)
// 	if err == nil {
// 		t.Fatalf("expected context cancelled error, received %v", err)
// 	}
// }

func TestGetProjectsBadJSON(t *testing.T) {
	// TODO: check that projects equals what we're expecting
	Setup()

	TestClient.BaseURL = TestServer.URL + "/bad-json"
	_, err := TestClient.Projects.GetProjects()
	if err == nil {
		t.Fatal("expected error, received nil")
	}
}

func TestGetProjectByID(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL
	_, err := TestClient.Projects.GetProjectByID(1)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestGetProjectByIDErrGetProjects(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/bad-json"
	_, err := TestClient.Projects.GetProjectByID(1)
	if err == nil {
		t.Fatalf("expected error, received nil")
	}
}

func TestGetProjectByName(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestGetProjectByNameErrGetProjects(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/bad-json"
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err == nil {
		t.Fatalf("expected error, received nil")
	}
}

func TestGetProjectByNameNoProjects(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/no-projects"
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestGetProjectByNameNotFound(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/not-found"
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestGetProjectByIDNoProjects(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/no-projects"
	_, err := TestClient.Projects.GetProjectByID(1)
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestGetProjectByIDNotFound(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL + "/not-found"
	_, err := TestClient.Projects.GetProjectByID(100)
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestCreateProject(t *testing.T) {
	Setup()

	TestClient.BaseURL = TestServer.URL
	newProject := &types.NewProject{
		Name:       "New Project",
		Color:      1,
		ParentID:   1,
		ChildOrder: 1,
		IsFavorite: 0,
	}

	err := TestClient.Projects.CreateProject(newProject)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestCreateProjectBadRequest(t *testing.T) {
	Setup()

	TestClient.BaseURL = "\t"
	newProject := &types.NewProject{
		Name:       "New Project",
		Color:      1,
		ParentID:   1,
		ChildOrder: 1,
		IsFavorite: 0,
	}

	err := TestClient.Projects.CreateProject(newProject)
	if err == nil {
		t.Fatal("expected error, received nil")
	}
}

// TODO: add ability to add context into individual resource requests
// func TestCreateProjectDoRequestContextCancel(t *testing.T) {
// 	Setup()
//
// 	TestClient.BaseURL = TestServer.URL
// 	d := time.Now().Add(1 * time.Second)
// 	ctx, cancel := context.WithDeadline(context.Background(), d)
// 	cancel()
//
// 	newProject := &types.NewProject{
// 		Name:       "New Project",
// 		Color:      1,
// 		ParentID:   1,
// 		ChildOrder: 1,
// 		IsFavorite: 0,
// 	}
//
// 	err := TestClient.Projects.CreateProject(newProject, ctx)
// 	if err == nil {
// 		t.Fatalf("expected context cancelled error, received %v", err)
// 	}
// }
