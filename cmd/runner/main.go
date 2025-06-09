package main

import (
	"context"
	"github.com/x3a-tech/configo"
	"github.com/x3a-tech/envo"
	"github.com/x3a-tech/logit-go"
	"i-o-bouns-tasks-api/internal/api"
	"i-o-bouns-tasks-api/internal/app"
	"i-o-bouns-tasks-api/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.run.main"

	cfg, _ := configo.MustLoad[config.Config]()
	env, err := envo.New(cfg.Env())
	logParams := logit.Params{
		AppConf:    &cfg.App,
		LoggerConf: &cfg.Logger,
		Env:        (*configo.Env)(env),
	}

	logger := logit.MustNewLogger(&logParams)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}
	ctx := logger.NewCtx(context.Background(), op, nil)

	logger.Info(ctx, "Сервис запущен успешно")

	tasksParams := api.TasksParams{
		Logger: logger,
	}

	tasks := api.NewTasks(&tasksParams)

	routerParams := api.Params{
		Task:   tasks,
		Logger: logger,
	}

	newApi := api.NewApi(&routerParams)

	appParams := app.Params{
		Api:    newApi,
		Logger: logger,
	}

	newApp := app.NewApp(&appParams)

	err = newApp.Init(ctx)
	if err != nil {
		logger.Error(ctx, err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-quit:
		logger.Info(ctx, "Завершение работы сервиса")
	}

	logger.Info(ctx, "Сервис успешно завершил работу")
}
