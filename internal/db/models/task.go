package models

import "time"

type Task struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":         t.Name,
		"description":  t.Description,
		"is_completed": t.IsCompleted,
	}
}
