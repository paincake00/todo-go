package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paincake00/todo-go/internal/db/postgres"
	httptodo "github.com/paincake00/todo-go/internal/delivery/http"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	config      Config
	logger      *zap.SugaredLogger
	taskHandler *httptodo.TaskHandler
	router      *gin.Engine
	gormDB      *gorm.DB
}

func NewApp(config Config, logger *zap.SugaredLogger) *App {
	app := &App{}

	app.config = config
	app.logger = logger

	gormDB, err := postgres.ConnectDB(
		config.db.address,
		config.db.maxOpenConn,
		config.db.maxIdleConn,
		config.db.maxConnLifetime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	app.gormDB = gormDB

	app.taskHandler = httptodo.NewTaskHandler(logger)

	app.router = app.InitRouter()

	return app
}

func (app *App) Run() error {
	// Канал для ошибок во время работы
	errChan := make(chan error, 10)

	// HTTP Server
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: app.router,
	}

	// graceful shutdown
	shutdown := make(chan error, 1) // канал для ошибок во время остановки
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Infof("Got signal %s, exiting gracefully...", s)

		// Context for shutdown of HTTP server
		ctxForSrv, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctxForSrv)
	}()

	// Starting HTTP server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("http server error: %w", err)
		}
	}()

	app.logger.Infow("Server started", "addr", srv.Addr)

	select {
	case err := <-errChan:
		app.logger.Errorf("fatal background error: %v", err)
		return err
	case err := <-shutdown:
		if err != nil {
			app.logger.Errorf("error occured on server shutting down: %s", err.Error())
			return err
		}
	}

	app.logger.Infof("Service is shutting down...")

	// sql db connection closing
	sqlDB, err := app.gormDB.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Close(); err != nil {
		app.logger.Errorf("error occured on closing database connection: %s", err.Error())
		return err
	}

	return nil
}
