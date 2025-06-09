package app

import (
	"context"
	"fmt"
	"github.com/x3a-tech/logit-go"
	"i-o-bouns-tasks-api/internal/api"
)

type App struct {
	api    *api.Api
	logger logit.Logger
}

type Params struct {
	Api    *api.Api
	Logger logit.Logger
}

func NewApp(params *Params) *App {
	return &App{
		api:    params.Api,
		logger: params.Logger,
	}
}

func (a *App) Init(ctx context.Context) error {
	const op = "app.Init"
	ctx = a.logger.NewOpCtx(ctx, op)

	err := a.api.Init(ctx)
	if err != nil {
		a.logger.Error(ctx, fmt.Errorf("ошибка при инициализации API: %v", err))
		return err
	}

	return nil
}
