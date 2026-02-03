package postgres

import (
	"context"
	"errors"

	"github.com/paincake00/todo-go/internal/db/models"
	"gorm.io/gorm"
)

var ErrorNotFound = errors.New("not found")

type ITaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Task, error)
	GetById(ctx context.Context, id uint) (*models.Task, error)
	UpdateById(ctx context.Context, task map[string]interface{}) (*models.Task, error)
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

func (t TaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	err := gorm.G[models.Task](t.db).Create(ctx, task)
	if err != nil {
		return nil, err
	}
	return t.GetById(ctx, task.Id)
}

func (t TaskRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Task, error) {
	tasks, err := gorm.G[models.Task](t.db).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskRepository) GetById(ctx context.Context, id uint) (*models.Task, error) {
	task, err := gorm.G[models.Task](t.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, ErrorNotFound
	}
	return &task, nil
}

func (t TaskRepository) UpdateById(ctx context.Context, task map[string]interface{}) (*models.Task, error) {
	//task.UpdatedAt = time.Now()

	i, _ := task["id"].(uint)

	t.db.Model(&models.Task{}).Where("id = ?", i).Updates(task)
	return t.GetById(ctx, i)
}

func (t TaskRepository) DeleteById(ctx context.Context, id uint) error {
	rowsAffected, err := gorm.G[models.Task](t.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorNotFound
	}
	return nil
}
