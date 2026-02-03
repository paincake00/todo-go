package http

import (
	"errors"
	httpgo "net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paincake00/todo-go/internal/db/postgres"
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

// CreateTask godoc
// @Summary      Создать задачу
// @Description  Создать новую задачу
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body dto.TaskCreateDTO true "Тело новой задачи"
// @Success      200  {object}  dto.TaskResponseDTO
// @Failure      400  {object}  dto.BadReqErrorResponse
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task dto.TaskCreateDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(httpgo.StatusBadRequest, dto.BadReqErrorResponse{Code: httpgo.StatusBadRequest, Error: err.Error()})
		return
	}

	ctx := c.Request.Context()

	taskModel, err := h.taskService.Create(ctx, mapper.FromTaskCreateDTOtoModel(&task))
	if err != nil {
		c.JSON(
			httpgo.StatusInternalServerError,
			dto.InternalErrorResponse{Code: httpgo.StatusInternalServerError, Error: err.Error()},
		)
		return
	}

	c.JSON(httpgo.StatusOK, mapper.FromTaskModelToDTO(taskModel))
}

// GetAllTasks godoc
// @Summary Получить список задач
// @Description Возвращает список задач с пагинацией
// @Tags tasks
// @Produce json
// @Param limit query int false "Limit (pagination)" default(20)
// @Param offset query int false "Offset (pagination)" default(0)
// @Success 200 {array} dto.TaskResponseDTO
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(LIMIT)))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(OFFSET)))

	tasks, err := h.taskService.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(
			httpgo.StatusInternalServerError,
			dto.InternalErrorResponse{Code: httpgo.StatusInternalServerError, Error: err.Error()},
		)
		return
	}
	c.JSON(httpgo.StatusOK, mapper.FromTaskModelListToDTO(tasks))
}

// GetTaskById godoc
// @Summary Получить задачу по ID
// @Description Возвращает одну задачу по идентификатору
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} dto.TaskResponseDTO
// @Failure      400  {object}  dto.BadReqErrorResponse
// @Failure      404  {object}  dto.NotFoundErrorResponse
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpgo.StatusBadRequest, dto.BadReqErrorResponse{Code: httpgo.StatusBadRequest, Error: err.Error()})
		return
	}

	task, err := h.taskService.GetById(c.Request.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrorNotFound):
			c.JSON(httpgo.StatusNotFound, dto.NotFoundErrorResponse{
				Code:  httpgo.StatusNotFound,
				Error: err.Error(),
			})
		default:
			c.JSON(httpgo.StatusInternalServerError, dto.InternalErrorResponse{
				Code:  httpgo.StatusInternalServerError,
				Error: err.Error(),
			})
		}
		return
	}
	c.JSON(httpgo.StatusOK, mapper.FromTaskModelToDTO(task))
}

// UpdateTaskById godoc
// @Summary Обновить задачу
// @Description Обновляет задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body dto.TaskResponseDTO true "Обновляемые поля"
// @Success 200 {object} dto.TaskResponseDTO
// @Failure      400  {object}  dto.BadReqErrorResponse
// @Failure      404  {object}  dto.NotFoundErrorResponse
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTaskById(c *gin.Context) {
	var task dto.TaskUpdateDTO // использовать map TODO

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(httpgo.StatusBadRequest, dto.BadReqErrorResponse{Code: httpgo.StatusBadRequest, Error: err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpgo.StatusBadRequest, dto.BadReqErrorResponse{Code: httpgo.StatusBadRequest, Error: err.Error()})
		return
	}
	task.Id = uint(id)

	if task.Name == nil && task.Description == nil && task.IsCompleted == nil {
		c.JSON(httpgo.StatusBadRequest, dto.BadReqErrorResponse{Code: httpgo.StatusBadRequest, Error: "at least one field must be provided"})
		return
	}

	taskModel, err := h.taskService.UpdateById(c.Request.Context(), mapper.FromTaskUpdateDTOtoMap(&task))
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrorNotFound):
			c.JSON(httpgo.StatusNotFound, dto.NotFoundErrorResponse{
				Code:  httpgo.StatusNotFound,
				Error: err.Error(),
			})
		default:
			c.JSON(httpgo.StatusInternalServerError, dto.InternalErrorResponse{
				Code:  httpgo.StatusInternalServerError,
				Error: err.Error(),
			})
		}
		return
	}

	c.JSON(httpgo.StatusOK, mapper.FromTaskModelToDTO(taskModel))
}

// DeleteTaskById godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по ID
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 204 "Задача удалена"
// @Failure      400  {object}  dto.BadReqErrorResponse
// @Failure      404  {object}  dto.NotFoundErrorResponse
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpgo.StatusBadRequest, dto.BadReqErrorResponse{Code: httpgo.StatusBadRequest, Error: err.Error()})
		return
	}

	err = h.taskService.DeleteById(c.Request.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrorNotFound):
			c.JSON(httpgo.StatusNotFound, dto.NotFoundErrorResponse{
				Code:  httpgo.StatusNotFound,
				Error: err.Error(),
			})
		default:
			c.JSON(httpgo.StatusInternalServerError, dto.InternalErrorResponse{
				Code:  httpgo.StatusInternalServerError,
				Error: err.Error(),
			})
		}
		return
	}

	c.JSON(httpgo.StatusNoContent, nil)
}
