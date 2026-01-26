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
	}
}
