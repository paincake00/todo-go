package service

import (
	"context"

	"github.com/paincake00/todo-go/internal/db/models"
	"github.com/paincake00/todo-go/internal/db/postgres"
)

type ITaskService interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Task, error)
	GetById(ctx context.Context, id uint) (*models.Task, error)
	UpdateById(ctx context.Context, task *models.Task) (*models.Task, error)
	DeleteById(ctx context.Context, id uint) error
}

type TaskService struct {
	taskRepository postgres.ITaskRepository
}

func NewTaskService(taskRepository postgres.ITaskRepository) ITaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (s *TaskService) Create(ctx context.Context, data *models.Task) (*models.Task, error) {
	return s.taskRepository.Create(ctx, data)
}

func (s *TaskService) GetAll(ctx context.Context, limit, offset int) ([]models.Task, error) {
	return s.taskRepository.GetAll(ctx, limit, offset)
}

func (s *TaskService) GetById(ctx context.Context, id uint) (*models.Task, error) {
	return s.taskRepository.GetById(ctx, id)
}

func (s *TaskService) UpdateById(ctx context.Context, task *models.Task) (*models.Task, error) {
	return s.taskRepository.UpdateById(ctx, task)
}

func (s *TaskService) DeleteById(ctx context.Context, id uint) error {
	return s.taskRepository.DeleteById(ctx, id)
}
