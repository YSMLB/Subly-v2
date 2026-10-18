package postgres

import (
	"context"
	"errors"

	"subly-v2/internal/domain"

	"github.com/jackc/pgx/v5"
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

func (ct *TaskRepository) GetTaskID (id int) (*domain.Task, error){
	query := "SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE id = $1"
	ctx := context.Background()

	var t domain.Task

	err := ct.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.UserID, &t.Title, 
		&t.Description, &t.Status, &t.Priority, 
		&t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows) == true{
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}