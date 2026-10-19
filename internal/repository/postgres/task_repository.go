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

func (gl *TaskRepository) GetList (userID int) ([]domain.Task, error){
	ctx := context.Background()
	query := "SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE user_id = $1"


	tSlice := make([]domain.Task, 0)

	rows, err := gl.pool.Query(ctx, query, userID)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var t domain.Task
		
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, 
			&t.Description, &t.Status, &t.Priority, 
			&t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		)

		if err != nil{
			return nil, err
		}

		tSlice = append(tSlice, t)
	}

	if err = rows.Err(); err != nil{
		return nil, err		
	}
	
	return tSlice, nil
}

func (ut *TaskRepository) UpdateTask (task *domain.Task) error{
	ctx := context.Background()
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = NOW() WHERE id = $6"

	_, err := ut.pool.Exec(ctx, query, task.Title, task.Description, 
		task. Status, task.Priority, 
		task.DueDate, task.ID)
	if err != nil{
		return err
	}

	return nil
}

func (dt *TaskRepository) DeleteTask (id int) error{
	ctx := context.Background()
	query := "DELETE FROM tasks WHERE id = $1"

	_, err := dt.pool.Exec(ctx, query, id)
	
	if err != nil{
		return err
	}

	return nil
}