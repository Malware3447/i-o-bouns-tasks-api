package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/x3a-tech/logit-go"
	"i-o-bouns-tasks-api/internal/service/tasks"
	"net/http"
)

type Api struct {
	router  *chi.Mux
	service *tasks.Service
	logger  logit.Logger
}

type Params struct {
	Service *tasks.Service
	Logger  logit.Logger
}

func NewApi(params *Params) *Api {
	return &Api{
		router:  nil,
		service: params.Service,
		logger:  params.Logger,
	}
}

func (a *Api) Init(ctx context.Context) error {
	const op = "app.Init.api.Init"
	ctx = a.logger.NewOpCtx(ctx, op)

	a.router = chi.NewRouter()

	a.router.Route("/api/v1", func(r chi.Router) {
		r.Post("/create_task", a.service.CreateTask)
		r.Get("/get_task/{id}", a.service.GetTask)
		r.Delete("/delete_task/{id}", a.service.DeleteTask)
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", 8080), a.router); err != nil {
			a.logger.Error(ctx, fmt.Errorf("ошибка запуска HTTP сервера: %w", err))
		}
	}()

	return nil
}
