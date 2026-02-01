package http

import (
	httpgo "net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paincake00/todo-go/internal/delivery/dto"
	"github.com/paincake00/todo-go/internal/domain/mapper"
	"github.com/paincake00/todo-go/internal/domain/service"
	"go.uber.org/zap"
)

const (
	LIMIT  = 10
	OFFSET = 0
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
		return
	}

	ctx := c.Request.Context()

	taskModel, err := h.taskService.Create(ctx, mapper.FromTaskDTOtoModel(&task))
	if err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(httpgo.StatusOK, mapper.FromTaskModelToDTO(taskModel))
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(LIMIT)))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(OFFSET)))

	tasks, err := h.taskService.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(httpgo.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpgo.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.GetById(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(httpgo.StatusOK, task)
}

func (h *TaskHandler) UpdateTaskById(c *gin.Context) {
	var task dto.TaskDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(httpgo.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpgo.StatusBadRequest, gin.H{"error": err.Error()})
	}
	task.Id = uint(id)

	taskModel, err := h.taskService.UpdateById(c.Request.Context(), mapper.FromTaskDTOtoModel(&task))
	if err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(httpgo.StatusOK, mapper.FromTaskModelToDTO(taskModel))
}

func (h *TaskHandler) DeleteTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpgo.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.taskService.DeleteById(c.Request.Context(), uint(id)); err != nil {
		c.JSON(httpgo.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
