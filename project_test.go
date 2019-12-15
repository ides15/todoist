package main

import (
	"testing"
)

func TestGetProjects(t *testing.T) {
	// TODO: check that projects equals what we're expecting
	setup()

	client.baseURL = testServer.URL
	projects, err := client.Projects.GetProjects()
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}

	expectedLength := 2
	if len(*projects) != expectedLength {
		t.Fatalf("expected %d projects, received %d", expectedLength, len(*projects))
	}
}

func TestGetProjectByID(t *testing.T) {
	setup()

	client.baseURL = testServer.URL
	_, err := client.Projects.GetProjectByID(1)
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}

func TestGetProjectByName(t *testing.T) {
	setup()

	client.baseURL = testServer.URL
	_, err := client.Projects.GetProjectByName("Inbox")
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
}
