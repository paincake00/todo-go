package mapper

import (
	"github.com/paincake00/todo-go/internal/db/models"
	"github.com/paincake00/todo-go/internal/delivery/dto"
)

func FromTaskDTOtoModel(taskDTO *dto.TaskDTO) *models.Task {
	return &models.Task{
		Id:          taskDTO.Id,
		Name:        taskDTO.Name,
		Description: taskDTO.Description,
		IsCompleted: taskDTO.IsCompleted,
		CreatedAt:   taskDTO.CreatedAt,
		UpdatedAt:   taskDTO.UpdatedAt,
	}
}

func FromTaskModelToDTO(task *models.Task) *dto.TaskDTO {
	return &dto.TaskDTO{
		Id:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
