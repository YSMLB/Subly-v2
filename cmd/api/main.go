package main

import (
	"fmt"
	"log"
	"net/http"
	"subly-v2/internal/domain"
	"subly-v2/internal/handler"
	"subly-v2/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// MockTaskRepo — фейковое хранилище для проверки сервиса без реальной базы данных
type MockTaskRepo struct{}

//MOckScheduleRepo - фейк хранилище для проверки сервиса без реальной бд
type MockScheduleRepo struct{}

//MockChangeRepo - фейк хранилизе ждя проверки сервиса без реальной бд
type MockChangeRepo struct{}

// Реализуем метод создания задачи (тот самый, который вызывает сервис)
func (m *MockTaskRepo) CreateTask(task domain.Task) error {
	fmt.Println("[MOCK DB] Задача успешно записана в базу:", task.Title)
	return nil
}

// Заглушки для остальных методов интерфейса TaskRepository
func (m *MockTaskRepo) GetTaskID(id int) (*domain.Task, error) {
	fmt.Println("[MOCK DB] Репозиторий подтвердил поиск задачи")
	return &domain.Task{
		ID: id,
		UserID: 3,
		Title: "title",
	}, nil
}
func (m *MockTaskRepo) GetList(userID int) ([]domain.Task, error)    {
	fmt.Println("[MOCK DB] Репозиторий подвердил поиск задач")
	return []domain.Task{
		{ID: userID,
		UserID: 6,
		Title: "саня ананьев чмо",
		Description: "углерод",},
		{ID: userID,
		UserID: 7,
		Title: "максименков тоже",
		Description: "азот",},
	}, nil
}
func (m *MockTaskRepo) UpdateTask(task *domain.Task) error {
	fmt.Println("[MOCK DB] Репозиторий подтвердил обновление задачи")
	return nil
}
func (m *MockTaskRepo) DeleteTask(taskID int) error {
	fmt.Println("[MOCK DB] Репозиторий подтвердил удаление задачи")
	return nil
}

func (m *MockScheduleRepo) AddSlot(schedule domain.ScheduleSlot) error{
	fmt.Println("[MOCK DB] Репозиторий подтвердил добавления занятий")
	return nil
}
func (m *MockScheduleRepo)GetScheduleToday(userID int, dayOfweek int) (*domain.ScheduleSlot , error){
	fmt.Println("[MOCK DB] Репозиторий подтверди поиск и выдачу занятий")
	return &domain.ScheduleSlot{
		UserID: 3,
		Subject: "МДК 05.02",
		Room: "401",
		StartTime: "08:30",
		EndTime: "10:00",
		DayOfWeek: 6,
		Type: domain.SlotTypeLecture,
		Parity: domain.WeekParityBoth,
	}, nil
}
func (m *MockScheduleRepo)UpdateTaskSchedule(schedule domain.ScheduleSlot) error{
	fmt.Println("[MOCK DB] Репозиторий подтвердил обновление занятий")
	return nil
}


//changes

func (m *MockChangeRepo)SaveChanges(changes domain.ProposedChange) (*domain.ProposedChange, error){
	fmt.Println("[MOCK DB] Репозиторий подтвердил сохранение изменений")
	return &domain.ProposedChange{
		ID: 1,
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "МДК 05.01",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "201",
		CreatedAdd: time.Now(),
	}, nil
}

func (m *MockChangeRepo)GetPendingByUserID(userID int) ([]domain.ProposedChange, error){
	fmt.Println("[MOCK DB] Репозиторий подтвердил просмотр изменений")
	return []domain.ProposedChange{
		{ID: 1,
		UserID: 1,
		SlotID: 1,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "МДК 05.01",
		NewStartTime: "10:10",
		NewEndTime: "11:40",
		NewRoom: "201",
		CreatedAdd: time.Now(),
	},
	{	ID: 2,
		UserID: 1,
		SlotID: 2,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "МДК 05.02",
		NewStartTime: "11:10",
		NewEndTime: "12:40",
		NewRoom: "301",
		CreatedAdd: time.Now(),},
	}, nil
}

func (m *MockChangeRepo) GetChangeByID(id int) (*domain.ProposedChange, error){
	fmt.Println("[MOCK DB] Репозиторий подтвердил поиск изменений")
	return &domain.ProposedChange{
		ID: 3,
		UserID: 4,
		SlotID: 5,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "Химия",
		NewStartTime: "14:00",
		NewEndTime: "15:30",
		NewRoom: "37",
		CreatedAdd: time.Now(),
	},nil
}

func (m *MockChangeRepo)UpdateStatusChanges(slotID int, userID int, newStatus domain.ChangeStatus) (*domain.ProposedChange, error){
	fmt.Println("[MOCK DB] Репозиторий подтвердил обновление изменений")
	return &domain.ProposedChange{
			ID: 4,
			UserID: 5,
			SlotID: 6,
			Action: domain.ChangeActionAddSlot,
			Status: domain.ChangeStatusPending,
			RawText: "Биология",
			NewStartTime: "14:00",
			NewEndTime: "15:30",
			NewRoom: "25",
			CreatedAdd: time.Now(),
	}, nil
}


func main() {

	mockChangeRepo := &MockChangeRepo{}
	mockTaskRepo := &MockTaskRepo{}
	mockScheduleRepo := &MockScheduleRepo{}

	changeService := &service.ChangeService{Repo: mockChangeRepo}
	taskService := &service.TaskService{Repo: mockTaskRepo}
	scheduleService := &service.ScheduleService{Repo: mockScheduleRepo}

	changeHandler := &handler.ChangeHandler{ChangeService: changeService}
	taskHandler := &handler.TaskHandler{ServiceTask: taskService}
	scheduleHandler := &handler.ScheduleHandler{ServiceSchedule: scheduleService}

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)

	//task
	router.Post("/tasks", taskHandler.CreateTask)
	router.Get("/tasks/{id}", taskHandler.GetTaskID)
	router.Get("/tasks", taskHandler.GetList)
	router.Patch("/tasks/{id}/status", taskHandler.UpdateStatus)
	router.Delete("/tasks/{id}", taskHandler.DeleteTask)

	//schedule
	router.Post("/schedule", scheduleHandler.AddSlot)
	router.Get("/schedule/today", scheduleHandler.GetScheduleToday)
	router.Put("/schedule/{id}", scheduleHandler.UpdateTaskSchedule)

	//change
	router.Post("/changes", changeHandler.SaveChanges)
	router.Get("/changes/pending", changeHandler.GetPendingByUserID)
	router.Get("/changes/{id}", changeHandler.GetChangeByID)
	router.Patch("/changes/{id}/status", changeHandler.UpdateStatusChanges)

	log.Println("Subly-v2 сервер запущен на порту:8080...")
	http.ListenAndServe(":8080", router)

}
