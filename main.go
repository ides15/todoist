package main

import (
	"context"
	"fmt"
)

func main() {
	todoist := Client{
		Token: "",
		// Instead of having to pass in a sync token each time, have a global
		// setting here to use sync tokens or just always sync fully
	}

	// GET ALL PROJECTS
	projects, err := todoist.GetProjects(context.Background(), "*")
	if err != nil {
		panic(err)
	}

	fmt.Println(projects)

	// GET PROJECT BY ID
	project, err := todoist.GetProjectByID(context.Background(), 2061811369, "*")
	if err != nil {
		panic(err)
	}

	fmt.Println(project)

	// GET PROJECT BY NAME
	project, err = todoist.GetProjectByName(context.Background(), "Inbox", "*")
	if err != nil {
		panic(err)
	}

	fmt.Println(project)

	// ADD PROJECT
	newProject := NewProject{
		Name:       "this is a new project",
		ParentID:   0,
		ChildOrder: 0,
		Color:      1,
		IsFavorite: 0,
	}

	new, err := todoist.AddProject(context.Background(), &newProject, nil, "*")
	if err != nil {
		panic(err)
	}

	fmt.Println(new)

	// UPDATE PROJECT
	updatedProject := UpdatedProject{
		ID:         new.ID,
		Name:       "this is an updated project",
		Color:      2,
		IsFavorite: 0,
		Collapsed:  0,
	}

	updated, err := todoist.UpdateProject(context.Background(), &updatedProject, "*")
	if err != nil {
		panic(err)
	}

	// MOVE PROJECT
	fmt.Println("Before move:", updated.ParentID)

	inbox, err := todoist.GetProjectByName(context.Background(), "Inbox", "*")
	if err != nil {
		panic(err)
	}

	moveProject := MovedProject{
		ID:       updated.ID,
		ParentID: inbox.ID,
	}

	err = todoist.MoveProject(context.Background(), &moveProject, "*")
	if err != nil {
		panic(err)
	}

	movedProject, err := todoist.GetProjectByID(context.Background(), updated.ID, "*")
	if err != nil {
		panic(err)
	}

	fmt.Println("After move:", movedProject.ParentID)

	// REORDER PROJECT
	fmt.Println("Before reorder:", movedProject.ChildOrder)

	reorderProject := ReorderedProject{
		Projects: &[]Reorder{{
			ID:         updated.ID,
			ChildOrder: 3,
		}},
	}

	err = todoist.ReorderProject(context.Background(), &reorderProject, "*")
	if err != nil {
		panic(err)
	}

	project, err = todoist.GetProjectByID(context.Background(), updated.ID, "*")
	if err != nil {
		panic(err)
	}

	fmt.Println("After reorder:", project.ChildOrder)

	// DELETE PROJECT
	deleteProject := DeletedProject{
		ID: updated.ID,
	}

	err = todoist.DeleteProject(context.Background(), &deleteProject, "*")
	if err != nil {
		panic(err)
	}
}
