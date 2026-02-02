package dto

import "time"

type TaskDTO struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	IsCompleted bool      `json:"is_completed" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
