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

func TestChangeRepository_GetPendingByUserID(t *testing.T){
	pool := setupTestDBChanges(t)
	defer pool.Close()
	pool1 := setupTestDBSchedule(t)
	defer pool1.Close()

	repo := NewChangeRepository(pool)
	repo1 := NewScheduleRepository(pool1)

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

	testChangeSlot1 := domain.ProposedChange{
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "test1 user1",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "201",
		CreatedAdd: time.Now(),
	}

	testChangeSlot2 := domain.ProposedChange{
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "test2 user1",
		NewStartTime: "08:30",
		NewEndTime: "10:00",
		NewRoom: "202",
		CreatedAdd: time.Now(),
	}

	testChangeSlot3 := domain.ProposedChange{
		UserID: 2,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "test1 user2",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "203",
		CreatedAdd: time.Now(),
	}

	testChangeSlot4 := domain.ProposedChange{
		UserID: 2,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "test2 user2",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "204",
		CreatedAdd: time.Now(),
	}

	testChangeSlot5 := domain.ProposedChange{
		UserID: 2,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusApproved,
		RawText: "test2 user2",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "205",
		CreatedAdd: time.Now(),
	}

	err := repo1.AddSlot(testScheduleSlot)
	if err != nil{
		t.Fatalf("AddSlot вернул ошибку: %v", err)
	}

	_, err = repo.SaveChanges(testChangeSlot1)
	if err != nil{
		t.Fatalf("testSave1 вернул ошибку: %v", err)
	}

	_, err = repo.SaveChanges(testChangeSlot2)
	if err != nil{
		t.Fatalf("testSave2 вренул ошибку: %v", err)
	}

	_, err = repo.SaveChanges(testChangeSlot3)
	if err != nil{
		t.Fatalf("testSave3 вернул ошибку: %v", err)
	}
	
	_, err = repo.SaveChanges(testChangeSlot4)
	if err != nil{
		t.Fatalf("testSave4 вернул ошибку: %v", err)
	}
	
	_, err = repo.SaveChanges(testChangeSlot5)
	if err != nil{
		t.Fatalf("testSave5 вернул ошибку: %v", err)
	}

	pendingList1, err := repo.GetPendingByUserID(1)
	if err != nil{
		t.Fatalf("GetPendingByUserID для пользователя 1 вернул ошибку: %v", err)
	}

	if len(pendingList1) != 2{
		t.Fatalf("Количество записей для pendingList1 не совпало с нужным. Должно быть 2, а получилось: %v", len(pendingList1))
	}

	pendingList2, err := repo.GetPendingByUserID(2)

	if len(pendingList2) != 2{
		t.Fatalf("Количество записей для pendingList2 не совпало с нужным. Должно быть 2, а получилось: %v", len(pendingList2))
	}

	if err != nil{
		t.Fatalf("GetPendingByUserID для пользователя 2 вернул ошибку: %v", err)
	}
}

