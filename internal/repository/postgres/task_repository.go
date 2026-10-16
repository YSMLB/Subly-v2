package postgres

import (
	"context"

	"subly-v2/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct{
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository{
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) CreateTask(task domain.Task) error{
	query := "INSERT INTO tasks (user_id, title, description, status, priority, due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := r.pool.Exec(context.Background(), query, task.UserID, task.Title, 
	task.Description, task.Status, task.Priority,
	 task.DueDate ,task.CreatedAt, task.UpdatedAt)
	if err != nil{
		return err
	}
	
	return nil
}