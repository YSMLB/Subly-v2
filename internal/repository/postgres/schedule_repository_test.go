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

func TestScheduleRepository_GetScheduleToday(t *testing.T){
	ctx := context.Background()
	pool := setupTestDBSchedule(t)
	defer pool.Close()

	repo := NewScheduleRepository(pool)

	testSchedule := domain.ScheduleSlot{
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

	err := repo.AddSlot(testSchedule)
	if err != nil{
		t.Fatalf("AddSlot завершился с ошибкой: %v", err)
	}

	test, err := repo.GetScheduleToday(testSchedule.UserID, 
		testSchedule.DayOfWeek)

	if err != nil{
		t.Fatalf("GetScheduleToday завершился с ошибкой: %v", err)
	}

	var g domain.ScheduleSlot
	query := "SELECT user_id, subject, teacher, room, day_of_week, start_time, end_time, type, parity, is_cancelled FROM schedule_slots WHERE user_id = $1 AND day_of_week = $2"

	err = pool.QueryRow(ctx, query, test.UserID, test.DayOfWeek).Scan(
		&g.UserID, &g.Subject, &g.Teacher,
		&g.Room, &g.DayOfWeek, &g.StartTime,
		&g.EndTime, &g.Type, &g.Parity,
		&g.IsCancelled,
	)

	if err != nil{
		t.Fatalf("Ошибка выдачи расписания в БД: %v", err)
	}
}

