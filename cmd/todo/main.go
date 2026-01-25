package main

import (
	"github.com/joho/godotenv"
	"github.com/paincake00/todo-go/internal/app"
	"github.com/paincake00/todo-go/internal/utils/logs"
)

func main() {
	logger := logs.NewLogger()
	defer func() {
		_ = logger.Sync()
	}()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	config := app.LoadConfig()

	application := app.NewApp(config, logger)

	err = application.Run()
	if err != nil {
		logger.Fatal(err)
	}
}
