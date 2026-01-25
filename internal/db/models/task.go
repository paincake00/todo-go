package models

import "time"

type Task struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
}
