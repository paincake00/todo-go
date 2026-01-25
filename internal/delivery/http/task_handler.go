package http

import (
	httpgo "net/http"

	"github.com/gin-gonic/gin"
	"github.com/paincake00/todo-go/internal/delivery/dto"
	"go.uber.org/zap"
)

type TaskHandler struct {
	logger *zap.SugaredLogger
}

func NewTaskHandler(logger *zap.SugaredLogger) *TaskHandler {
	return &TaskHandler{
		logger: logger,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task dto.TaskDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(httpgo.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {

}

func (h *TaskHandler) GetTaskById(c *gin.Context) {}

func (h *TaskHandler) UpdateTaskById(c *gin.Context) {}

func (h *TaskHandler) DeleteTaskById(c *gin.Context) {}
