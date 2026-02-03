package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paincake00/todo-go/internal/db/models"
	"github.com/paincake00/todo-go/internal/delivery/dto"
	"github.com/paincake00/todo-go/internal/domain/service"
	"github.com/paincake00/todo-go/internal/utils/logs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgrestodo "github.com/paincake00/todo-go/internal/db/postgres"
)

const (
	SCHEMA = "postgres"
	DbPass = "pass"
	DbUser = "user"
	DbName = "testdb"
)

var app *testApp

type testApp struct {
	gormDB         *gorm.DB
	taskRepository postgrestodo.ITaskRepository
	taskService    service.ITaskService
	taskHandler    *TaskHandler
	router         *gin.Engine
}

func setupTestDB(ctx context.Context) (testcontainers.Container, *gorm.DB) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.6",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": DbPass,
			"POSTGRES_USER":     DbUser,
			"POSTGRES_DB":       DbName,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("could not start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("could not get host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("could not get port for container: %v", err)
	}

	uri := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		SCHEMA, DbUser, DbPass, host, port.Port(), DbName)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	return container, db
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, gormDB := setupTestDB(ctx)

	taskRepository := postgrestodo.NewTaskRepository(gormDB)

	taskService := service.NewTaskService(taskRepository)

	taskHandler := NewTaskHandler(logs.NewLogger(), taskService)

	//gin.SetMode(gin.TestMode)
	router := gin.Default()
	//router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetAllTasks)
			tasks.GET("/:id", taskHandler.GetTaskById)
			tasks.PUT("/:id", taskHandler.UpdateTaskById)
			tasks.DELETE("/:id", taskHandler.DeleteTaskById)
		}
	}

	// init test app
	app = &testApp{
		gormDB:         gormDB,
		taskRepository: taskRepository,
		taskService:    taskService,
		taskHandler:    taskHandler,
		router:         router,
	}

	exitCode := m.Run()

	// cleanup containers
	err := container.Terminate(ctx)
	if err != nil {
		log.Fatalf("could not terminate container: %v", err)
	}

	os.Exit(exitCode)
}

func TestCreateTask(t *testing.T) {
	body := `{
	    "name": "task 1",
		"description": "This is a task 1"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response dto.TaskResponseDTO
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	require.NotZero(t, response.Id)
	require.Equal(t, "task 1", response.Name)
	t.Log("dto: ", response)

	t.Run("Get model from db", func(t *testing.T) {
		taskModel, err := app.taskRepository.GetById(context.Background(), response.Id)
		if err != nil {
			t.Errorf("could not get task by id: %v", err)
		}

		t.Log("model: ", *taskModel)
	})

	t.Run("non-name, description", func(t *testing.T) {
		body = `{
			"description": "This is a task 1"
		}`

		req = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("name, non-description", func(t *testing.T) {
		body = `{
			"name": "task 2"
		}`

		req = httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var response1 dto.TaskResponseDTO
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response1))

		t.Log("dto: ", response1)
	})
}

func TestGetAllTasks(t *testing.T) {
	//newTask := &models.Task{
	//	Name:        "task 2",
	//	Description: "This is a task 2",
	//	IsCompleted: true,
	//}
	//
	//createdTask, err := app.taskRepository.Create(context.Background(), newTask)
	//require.NoError(t, err)
	//
	//t.Logf("CreatedAt=%v UpdatedAt=%v", createdTask.CreatedAt, createdTask.UpdatedAt)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.TaskResponseDTO
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	t.Log("all tasks (dto): ", response)

	require.NotZero(t, len(response))
}

func TestGetById(t *testing.T) {
	testId := uint(1)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%d", testId), nil)
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response dto.TaskResponseDTO
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	require.NotEqual(t, "0001-01-01 00:00:00 +0000 UTC", response.CreatedAt.String())

	t.Log("dto: ", response)

	t.Run("Get model from db", func(t *testing.T) {
		taskModel, err := app.taskRepository.GetById(context.Background(), testId)
		if err != nil {
			t.Errorf("could not get task by id: %v", err)
		}
		t.Log("model: ", *taskModel)
	})
}

func TestUpdateTaskById(t *testing.T) {
	body := `{
      "name": "New task",
	  "description": "This is a new task",
	  "is_completed": true
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/1", strings.NewReader(body))
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response dto.TaskResponseDTO
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	t.Log("update task with id 1 (dto): ", response)

	t.Run("non-name", func(t *testing.T) {
		// description как поле передается, должен измениться на пустую строку
		body = `{
		  "description": "",
		  "is_completed": false
		}`

		req = httptest.NewRequest(http.MethodPut, "/api/v1/tasks/1", strings.NewReader(body))
		w = httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var response1 dto.TaskResponseDTO
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response1))

		t.Log("dto: ", response1)
	})

	t.Run("empty body", func(t *testing.T) {
		// description как поле передается, должен измениться на пустую строку
		body = `{}`

		req = httptest.NewRequest(http.MethodPut, "/api/v1/tasks/1", strings.NewReader(body))
		w = httptest.NewRecorder()

		app.router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteTaskById(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}
