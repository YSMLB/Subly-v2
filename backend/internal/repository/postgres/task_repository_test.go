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

func TestTaskRepository_GetTaskID(t *testing.T){
	ctx := context.Background()
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)

	testTask := domain.Task{
		UserID: 1,
		ID: 1,
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
	test, err := repo.GetTaskID(testTask.ID)
	if err != nil{
		t.Fatalf("GetTaskID завершился с ошибкой %v", err)
	}

	var te domain.Task
	query := "SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE id = $1"

	err = pool.QueryRow(ctx, query, test.ID).Scan(&te.ID, &te.UserID, 
		&te.Title, 
		&te.Description, &te.Status, &te.Priority, 
		&te.DueDate, &te.CreatedAt, &te.UpdatedAt)

	if err != nil{
		t.Fatalf("ошибка проверки наличия задачи в БД: %v", err )
	}
}

func TestTaskRepository_GetList(t *testing.T){
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)

	testTask := domain.Task{
		UserID: 1,
		ID: 1,
		Title: "Сдать отчет по БД",
		Description: "Лаба номер 3",
		Status: domain.TaskStatusPending,
		Priority: domain.TaskPriorityHigh,
		DueDate: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testTask2 := domain.Task{
		UserID: 1,
		ID: 1,
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

	err = repo.CreateTask(testTask2)
	if err != nil{
		t.Fatalf("CreateTask завершился с ошибкой: %v", err)
	}

	tasks, err := repo.GetList(testTask.UserID)
	if err != nil{
		t.Fatalf("GetList завершился с ошибкой: %v", err)
	}	

	if len(tasks) != 2{
		t.Fatal("Не 2 задача")
	}
	
}

func TestTaskRepository_UpdateTask(t *testing.T){
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewTaskRepository(pool)
	
	testTask := domain.Task{
		UserID: 1,
		ID: 1,
		Title: "Сдать отчет по БД",
		Description: "Лаба номер 3",
		Status: domain.TaskStatusPending,
		Priority: domain.TaskPriorityHigh,
		DueDate: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testTask2 := domain.Task{
		UserID: 1,
		ID: 2,
		Title: "Сдать отчет по БД 2",
		Description: "Лаба номер 3 2",
		Status: domain.TaskStatusPending,
		Priority: domain.TaskPriorityHigh,
		DueDate: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testTaskUpdate := domain.Task{
		UserID: 1,
		ID: 1,
		Title: "Мой день рождения",
		Description: "сегодня(19 октября)",
		Status: domain.TaskStatusInProgress,
		Priority: domain.TaskPriorityLow,
		DueDate: time.Now().Add(24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateTask(testTask)
	if err != nil{
		t.Fatalf("CreateTask1 вернул ошибку: %v", err)
	}

	err = repo.CreateTask(testTask2)
	if err != nil{
		t.Fatalf("CreateTask2 вернул ошибку: %v", err)
	}

	err = repo.UpdateTask(&testTaskUpdate)
	if err != nil{
		t.Fatalf("UpdateTask вернул ошибку: %v", err)
	}
}

func TestTaskRepository_DeleteTask(t *testing.T){
	pool := setupTestDB(t)
	defer pool.Close()
	id := 2

	repo := NewTaskRepository(pool)

	testTask1 := domain.Task{
		ID: 3,
		UserID: 2,
		Title: "delete task",
		Description: "hb to me",
		Status: domain.TaskStatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
		Priority: domain.TaskPriorityHigh,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testTask2 := domain.Task{
		ID: 4,
		UserID: 2,
		Title: "delete task 2",
		Description: "hb to me 2",
		Status: domain.TaskStatusPending,
		Priority: domain.TaskPriorityHigh,
		DueDate: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateTask(testTask1)
	if err != nil{
		t.Fatalf("CreateTask1 вернул ошибку: %v", err)
	}

	err = repo.CreateTask(testTask2)
	if err != nil{
		t.Fatalf("CreateTask2 верунл ошибку: %v", err)
	}

	err = repo.DeleteTask(id)
	if err != nil{
		t.Fatalf("DeleteTask вернул ошибку: %v", err)
	}
}