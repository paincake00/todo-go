package models

import "time"

type Task struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	IsCompleted bool
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
