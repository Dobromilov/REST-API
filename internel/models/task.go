package models

import (
	"time"
)

type Task struct {
	ID          int       `json:"id" db:"id"` // тег как поле будет называться в json и тег для библиотеки в sqlx
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// без указателей здесь бы была беда, а именно я бы не мог обновлять только конкретное поле,
// т.к. при работе без них у меня бы создавалось значение по умолчанию
