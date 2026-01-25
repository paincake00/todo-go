package app

import "go.uber.org/zap"

type App struct {
	config Config
	logger *zap.SugaredLogger
}

func NewApp(config Config, logger *zap.SugaredLogger) *App {
	app := &App{}

	app.config = config
	app.logger = logger

	return app
}

func (app *App) Run() error {
	return nil
}
