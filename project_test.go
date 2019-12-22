package todoist

import (
	"testing"

	"github.com/ides15/todoist/types"
)

func TestGetProjects(t *testing.T) {
	// TODO: check that projects equals what we're expecting
	Setup()

	TestClient.baseURL = TestServer.URL
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

	TestClient.baseURL = "\t"
	_, err := TestClient.Projects.GetProjects()
	if err == nil {
		t.Fatalf("expected error, received nil")
	}
}

// TODO: add ability to add context into individual resource requests
// func TestGetProjectsDoRequestContextCancel(t *testing.T) {
// 	Setup()

// 	TestClient.baseURL = TestServer.URL
// 	d := time.Now().Add(1 * time.Second)
// 	ctx, cancel := context.WithDeadline(context.Background(), d)
// 	cancel()

// 	_, err := TestClient.Projects.GetProjects(ctx)
// 	if err == nil {
// 		t.Fatalf("expected context cancelled error, received %v", err)
// 	}
// }

func TestGetProjectsBadJSON(t *testing.T) {
	// TODO: check that projects equals what we're expecting
	Setup()

	TestClient.baseURL = TestServer.URL + "/bad-json"
	_, err := TestClient.Projects.GetProjects()
	if err == nil {
		t.Fatal("expected error, received nil")
	}
}

func TestGetProjectByID(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL
	_, err := TestClient.Projects.GetProjectByID(1)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestGetProjectByIDErrGetProjects(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL + "/bad-json"
	_, err := TestClient.Projects.GetProjectByID(1)
	if err == nil {
		t.Fatalf("expected error, received nil")
	}
}

func TestGetProjectByName(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestGetProjectByNameErrGetProjects(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL + "/bad-json"
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err == nil {
		t.Fatalf("expected error, received nil")
	}
}

func TestGetProjectByNameNoProjects(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL + "/no-projects"
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestGetProjectByNameNotFound(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL + "/not-found"
	_, err := TestClient.Projects.GetProjectByName("Inbox")
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestGetProjectByIDNoProjects(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL + "/no-projects"
	_, err := TestClient.Projects.GetProjectByID(1)
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}

func TestGetProjectByIDNotFound(t *testing.T) {
	Setup()

	TestClient.baseURL = TestServer.URL + "/not-found"
	_, err := TestClient.Projects.GetProjectByID(100)
	if err == nil {
		t.Fatal("expected error, received nil")
	} else if err != types.ErrNotFound {
		t.Fatalf("expected %v, received %v", types.ErrNotFound, err)
	}
}
