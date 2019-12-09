package main

import (
	"github.com/ides15/todoist/types"
)

type ProjectService struct {
	c *Client
}

func (s *ProjectService) GetProjects() (*[]types.Project, error) {
	return nil, nil
}
