package mapper

import (
	"github.com/paincake00/todo-go/internal/db/models"
	"github.com/paincake00/todo-go/internal/delivery/dto"
)

func FromTaskCreateDTOtoModel(taskDTO *dto.TaskCreateDTO) *models.Task {
	return &models.Task{
		Name:        taskDTO.Name,
		Description: taskDTO.Description,
	}
}

func FromTaskUpdateDTOtoMap(taskDTO *dto.TaskUpdateDTO) map[string]interface{} {
	updatedFields := make(map[string]interface{})

	updatedFields["id"] = taskDTO.Id

	if taskDTO.Name != nil {
		updatedFields["name"] = *taskDTO.Name
	}
	if taskDTO.Description != nil {
		updatedFields["description"] = *taskDTO.Description
	}
	if taskDTO.IsCompleted != nil {
		updatedFields["is_completed"] = *taskDTO.IsCompleted
	}
	return updatedFields
}

func FromTaskModelToDTO(task *models.Task) *dto.TaskResponseDTO {
	return &dto.TaskResponseDTO{
		Id:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func FromTaskModelListToDTO(tasks []models.Task) []*dto.TaskResponseDTO {
	result := make([]*dto.TaskResponseDTO, len(tasks)) // length == capacity
	for i, task := range tasks {
		result[i] = FromTaskModelToDTO(&task)
	}
	return result
}
