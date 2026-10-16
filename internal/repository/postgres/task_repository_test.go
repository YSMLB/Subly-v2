package postgres

import (
	"context"
	"subly-v2/internal/domain"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const testStr = "postgres://postgres:postgres@localhost:5432/subly_db?sslmode=disable"

func setupTestDB(t *testing.T) *pgxpool.Pool{
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, testStr)

	if err != nil{
		t.Fatalf("Не удалось подключитсья к БД %v", err)
	}

	_ , err = pool.Exec(ctx, "TRUNCATE TABLE tasks RESTART IDENTITY CASCADE;")
	if err != nil{
		t.Fatalf("Ошибка очистки таблицы tasks %v", err)
	} 

	return pool
}

func TestTaskRepository_CreateTask(t *testing.T){
	ctx := context.Background()
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)

	testTask := domain.Task{
		UserID: 1,
		Title: "Сдать отчет по БД",
		Description: "Лаба номер 3",
		Status: domain.TaskStatusPending,
		Priority: domain.TaskPriorityHigh,
		DueDate: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateTask(testTask)
	if err != nil{
		t.Fatalf("CreateTask завершился с ошибкой: %v", err)
	}

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND title = $2",
	 testTask.UserID, testTask.Title).Scan(&count)
	
	if err != nil{
		t.Fatalf("ошибка проверки наличия задачи в БД: %v", err )
	}

	if count != 1 {
		t.Errorf("ожидалась 1 задача в базе, но найдено: %d", err)
	}
}