package postgres

import (
	"context"

	"github.com/paincake00/todo-go/internal/db/models"
	"gorm.io/gorm"
)

type ITaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetAll(ctx context.Context, limit, offset int) ([]models.Task, error)
	GetById(ctx context.Context, id uint) (*models.Task, error)
	UpdateById(ctx context.Context, task *models.Task) error
	DeleteById(ctx context.Context, id uint) error
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (repo *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	return gorm.G[models.Task](repo.db).Create(ctx, task)
}

func (repo *TaskRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Task, error) {
	tasks, err := gorm.G[models.Task](repo.db).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *TaskRepository) GetById(ctx context.Context, id uint) (*models.Task, error) {
	panic("implement me")
}

func (repo *TaskRepository) UpdateById(ctx context.Context, task *models.Task) error {
	panic("implement me")
}

func (repo *TaskRepository) DeleteById(ctx context.Context, id uint) error {
	panic("implement me")
}
