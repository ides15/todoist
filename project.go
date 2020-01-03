package todoist

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ides15/todoist/types"
)

type ProjectService struct {
	c *Client
}

func (s *ProjectService) GetProjects() ([]*types.Project, *types.Response, error) {
	s.c.Log("GetProjects called")

	req, err := s.c.NewRequest("*", nil, &[]string{"projects"})
	if err != nil {
		return nil, nil, err
	}

	response := new(types.Response)
	_, err = s.c.Do(context.Background(), req, response)
	if err != nil {
		return nil, nil, err
	}

	return response.Projects, response, nil
}

func (s *ProjectService) GetProjectByID(id int) (*types.Project, *types.Response, error) {
	s.c.Log("GetProjectByID called")

	projects, res, err := s.GetProjects()
	if err != nil {
		return nil, res, err
	}

	s.c.Log("projects", projects)

	for _, project := range projects {
		if project.ID == id {
			return project, res, nil
		}

		return nil, res, types.ErrNotFound
	}

	return nil, res, types.ErrNotFound
}

func (s *ProjectService) GetProjectByName(name string) (*types.Project, *types.Response, error) {
	s.c.Log("GetProjectByName called")

	projects, res, err := s.GetProjects()
	if err != nil {
		return nil, res, err
	}

	for _, project := range projects {
		if project.Name == name {
			return project, res, nil
		}

		return nil, res, types.ErrNotFound
	}

	return nil, res, types.ErrNotFound
}

func (s *ProjectService) CreateProject(p *types.NewProject) (*types.Response, error) {
	s.c.Log("CreateProject called")

	commands := &[]types.Command{
		{
			Type:   "project_add",
			TempID: uuid.New().String(),
			UUID:   uuid.New().String(),
			Args:   p,
		},
	}
	commandsString, _ := json.Marshal(commands)
	s.c.Logf("\tCommands: %v\n", string(commandsString))

	req, err := s.c.NewRequest("*", commands, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(context.Background(), req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProjectService) UpdateProject(p *types.UpdatedProject) (*types.Response, error) {
	s.c.Log("UpdateProject called")

	commands := &[]types.Command{
		{
			Type:   "project_update",
			TempID: uuid.New().String(),
			UUID:   uuid.New().String(),
			Args:   p,
		},
	}
	commandsString, _ := json.Marshal(commands)
	s.c.Logf("\tCommands: %v\n", string(commandsString))

	req, err := s.c.NewRequest("*", commands, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(context.Background(), req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProjectService) MoveProject(ID int, ParentID int) (*types.Response, error) {
	s.c.Log("MoveProject called")

	args := make(map[string]int)
	args["id"] = ID
	args["parent_id"] = ParentID

	commands := &[]types.Command{
		{
			Type:   "project_move",
			TempID: uuid.New().String(),
			UUID:   uuid.New().String(),
			Args:   args,
		},
	}
	commandsString, _ := json.Marshal(commands)
	s.c.Logf("\tCommands: %v\n", string(commandsString))

	req, err := s.c.NewRequest("*", commands, &[]string{"projects"})
	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(context.Background(), req, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
