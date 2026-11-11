package service

import (
	"context"
	"subly-v2/internal/domain"
	"subly-v2/internal/repository/postgres"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const testStr = "postgres://postgres:postgres@localhost:5432/subly_db?sslmode=disable"

func setupTestDB(t *testing.T) *pgxpool.Pool{
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, testStr)
	query := "TRUNCATE TABLE proposed_changes, schedule_slots RESTART IDENTITY CASCADE;"

	if err != nil{
		t.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	_, err = pool.Exec(ctx, query)
	if err != nil{
		t.Fatalf("Ошибка очистки таблиц: %v", err)
	}

	return pool
}

func TestChangeService_UpdateStatusChanges_approved_Integration(t *testing.T){
	pool := setupTestDB(t)
	defer pool.Close()

	scheduleRepo := postgres.NewScheduleRepository(pool)
	changeRepo := postgres.NewChangeRepository(pool)	

	changeService := &ChangeService{Repo: changeRepo, RepoSchedule: scheduleRepo}

	TestAddSlot := domain.ScheduleSlot{
		ID: 1,
		UserID:      1,
		Subject:     "MDK 05.02",
		Teacher:     "Ivanova V. I.",
		Room:        "21",
		DayOfWeek:   2,
		StartTime:   "15:40",
		EndTime:     "17:10",
		Type:        domain.SlotTypePractice,
		Parity:      domain.WeekParityEven,
		IsCancelled: false,
	}

	TestSaveChanges := domain.ProposedChange{
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionChangeRoom,
		Status: domain.ChangeStatusPending,
		RawText: "English",
		NewStartTime: "08:30",
		NewEndTime: "10:00",
		NewRoom: "209",
		CreatedAdd: time.Now(),
	}
	

	err := scheduleRepo.AddSlot(TestAddSlot)
	if err != nil{
		t.Fatalf("AddSlot вернул ошибку: %v", err)
	}

	_, err = scheduleRepo.GetScheduleToday(TestAddSlot.UserID, TestAddSlot.DayOfWeek)
	if err != nil{
		t.Fatalf("GetScheduleToday вернул ошибку: %v", err)
	}

	save, err := changeRepo.SaveChanges(TestSaveChanges)
	if err != nil{
		t.Fatalf("SaveChanges вернул ошибку: %v", err)
	}

	testNewMethod, err := changeService.UpdateStatusChanges(save.ID, 1, domain.ChangeStatusApproved)

	if err !=nil{
		t.Fatalf("UpdateStatusChanges вернул ошибку: %v", err)
	}
	if testNewMethod.Status != domain.ChangeStatusApproved{
		t.Fatalf("Получили не тот статус: %v", testNewMethod.Status)
	}

	getSlot, err := scheduleRepo.GetSlotByID(testNewMethod.SlotID)
	if err != nil{
		t.Fatalf("GetSlotByID вернул ошибку: %v", err)
	}

	if getSlot.Room != "209"{
		t.Fatalf("Не тот кабинет: %v", getSlot.Room)
	}
}