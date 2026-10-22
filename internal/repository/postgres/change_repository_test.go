package postgres

import (
	"context"
	"subly-v2/internal/domain"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDBChanges(t *testing.T) *pgxpool.Pool{
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, testStr)
	query := "TRUNCATE TABLE proposed_changes RESTART IDENTITY CASCADE;"

	if err != nil{
		t.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	_, err = pool.Exec(ctx, query)
	if err != nil{
		t.Fatalf("Ошибка чистки таблицы proposed_changes: %v", err)
	}

	return pool
}

func TestChangeRepository_SaveChanges(t *testing.T){
	pool := setupTestDBChanges(t)
	pool1 := setupTestDBSchedule(t)
	defer pool.Close()

	repo := NewChangeRepository(pool)
	repo1 := NewScheduleRepository(pool1)

	testChangeSlot := domain.ProposedChange{
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "English",
		NewStartTime: "08:30",
		NewEndTime: "10:00",
		NewRoom: "209",
		CreatedAdd: time.Now(),
	}
	
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

	err := repo1.AddSlot(testScheduleSlot)
	if err != nil{
		t.Fatalf("Не удалось создать Тестовое расписание для SaveChanges: %v", err)
	}

	_, err = repo.SaveChanges(testChangeSlot)
	if err != nil{
		t.Fatalf("SaveChanges завершил работу с ошибкой: %v", err)
	}
}