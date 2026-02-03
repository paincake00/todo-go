package dto

import "time"

/*
	Надо для определнных эндпоинтов (на POST, PUT, PATCH и тд)
	создавать отдельные DTO с определнным набором полей,
    которые будут предоставляться фронтенду через API.
*/

type TaskCreateDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type TaskUpdateDTO struct {
	Id          uint    `json:"-"`    // не будет видно при ShouldBindJson, заполнится из path param вручную
	Name        *string `json:"name"` // указатели для того, чтобы отличать отствуие поля в переданном JSON от zero value
	Description *string `json:"description"`
	IsCompleted *bool   `json:"is_completed"`
}

type TaskResponseDTO struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
