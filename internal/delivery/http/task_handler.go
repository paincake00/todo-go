package http

import (
	httpgo "net/http"

	"github.com/gin-gonic/gin"
	"github.com/paincake00/todo-go/internal/delivery/dto"
	"github.com/paincake00/todo-go/internal/domain/mapper"
	"github.com/paincake00/todo-go/internal/domain/service"
	"go.uber.org/zap"
)

type TaskHandler struct {
	taskService service.ITaskService
	logger      *zap.SugaredLogger
}

func NewTaskHandler(logger *zap.SugaredLogger, taskService service.ITaskService) *TaskHandler {
	return &TaskHandler{
		logger:      logger,
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task dto.TaskDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(httpgo.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx := c.Request.Context()

	taskModel := mapper.FromTaskDTOtoModel(&task)

	if err := h.taskService.Create(ctx, taskModel); err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(httpgo.StatusOK, gin.H{"data": mapper.FromTaskModelToDTO(taskModel)})
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAll(c.Request.Context(), 10, 0)
	if err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(httpgo.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskById(c *gin.Context) {}

func (h *TaskHandler) UpdateTaskById(c *gin.Context) {}

func (h *TaskHandler) DeleteTaskById(c *gin.Context) {}
