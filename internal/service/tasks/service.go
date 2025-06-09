package tasks

import (
	"i-o-bouns-tasks-api/internal/api/tasks"
	"net/http"
)

type Service struct {
	repo tasks.TaskHandler
}

type Params struct {
	Repo tasks.TaskHandler
}

func NewService(params *Params) *Service {
	return &Service{
		repo: params.Repo,
	}
}

func (s *Service) CreateTask(w http.ResponseWriter, r *http.Request) {
	s.repo.CreateTask(w, r)
}

func (s *Service) GetTask(w http.ResponseWriter, r *http.Request) {
	s.repo.GetTask(w, r)
}

func (s *Service) DeleteTask(w http.ResponseWriter, r *http.Request) {
	s.repo.DeleteTask(w, r)
}
