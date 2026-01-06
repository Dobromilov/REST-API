package database

import (
	"database/sql"
	"errors"
	"simple-api/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskScore struct {
	db *sqlx.DB
}

func NewTaskScore(db *sqlx.DB) *TaskScore {
	return &TaskScore{db: db}
}

func (s *TaskScore) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	query := `
SELECT id, title, description, completed, created_at, updated_at 
FROM tasks 
order by created_at desc;`

	err := s.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskScore) GetByID(ID int) (*models.Task, error) {
	var task models.Task

	query := `
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
WHERE id = $1;` // плейсхолдер - заполнитель для параметра (против sql инъекции)

	err := s.db.Get(&task, query, ID)

	if err == sql.ErrNoRows {
		return &models.Task{}, nil
	} // нет в базе данных

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskScore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	query := `
INSERT INTO tasks (title, description, completed, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING  title, description, completed, created_at, updated_at;`

	now := time.Now()

	err := s.db.QueryRowx(query, input.Title, input.Description, input.Completed, now, now).StructScan(&task) // вставляем в наш query текущие значения; подставляем в структуру
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskScore) Update(ID int, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := s.GetByID(ID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	task.UpdatedAt = time.Now()

	query := `
UPDATE tasks 
SET title = $1, description = $2, completed = $3, updated_at = $4
WHERE id = $5
RETURNING id, title, description, completed, created_at, updated_at;`

	var updatedTask models.Task
	err = s.db.QueryRowx(query, task.Title, task.Description, task.Completed,
		task.UpdatedAt, task.ID).StructScan(&updatedTask)
	if err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

func (s *TaskScore) Delete(ID int) error {
	query := `DELETE FROM tasks WHERE id = $1;`

	result, err := s.db.Exec(query, ID) //выполняет запрос без возврата строк

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Task not found")
	}

	return nil
}
