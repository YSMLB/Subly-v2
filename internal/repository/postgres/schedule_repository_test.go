package postgres

import (
	"context"
	"subly-v2/internal/domain"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const testStrScheduler = "postgres://postgres:postgres@localhost:5432/subly_db?sslmode=disable"

func setupTestDBSchedule(t *testing.T) *pgxpool.Pool{
	ctx := context.Background()
	query := "TRUNCATE TABLE schedule_slots RESTART IDENTITY CASCADE;"
	pool, err := pgxpool.New(ctx, testStrScheduler)

	if err != nil{
		t.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	_, err = pool.Exec(ctx, query)
	if err != nil{
		t.Fatalf("ошибка очистки таблицы schedule_slots: %v", err)
	}

	return pool
}

func TestScheduleRepository_AddSlot(t *testing.T){
	pool := setupTestDBSchedule(t)
	defer pool.Close()

	repo := NewScheduleRepository(pool)

	testScheduleSlot := domain.ScheduleSlot{
		UserID: 1,
		Subject: "MDK 05.02",
		Teacher: "Ivanova V. I.",
		Room: "21",
		DayOfWeek: 2,
		StartTime: "15:40",
		EndTime: "17:10",
		Type: domain.SlotTypePractice,
		Parity: domain.WeekParityEven,
		IsCancelled: false,
	}

	err := repo.AddSlot(testScheduleSlot)
	if err != nil{
		t.Fatalf("AddSlot завершился с ошибкой: %v", err)
	}
}

