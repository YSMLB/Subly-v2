package main

import (
	"fmt"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
)

// MockTaskRepo — фейковое хранилище для проверки сервиса без реальной базы данных
type MockTaskRepo struct{}

//MOckScheduleRepo - фейк хранилище для проверки сервиса без реальной бд
type MockScheduleRepo struct{}

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
func (m *MockScheduleRepo)GetScheduleToday() error{
	return nil
}
func (m *MockScheduleRepo)UpdateTaskSchedule(){}
func main() {
	fmt.Println("--- Запуск проверки логики Subly ---")

	// 1. Создаем заглушку репозитория
	mockRepo := &MockTaskRepo{}

	// 1. Создали заглушку репы
	mockRepoSchedule := &MockScheduleRepo{}

	// 2. Инициализируем наш сервис задач
	taskService := &service.TaskService{
		Repo: mockRepo,
	}

	scheduleService := &service.ScheduleService{
		Repo: mockRepoSchedule,
	}

	// -------------------------------------------------------------
	// ТЕСТ 1: Попытка создать нормальную задачу
	// -------------------------------------------------------------
	fmt.Println("\n[ТЕСТ 1] Создаем валидную задачу...")
	goodTask := domain.Task{
		UserID:      1,
		Title:       "Сдать лабу по Go",
		Description: "Разобраться с микросервисом Subly",
	}

	createdTask, err := taskService.CreateTask(goodTask)
	if err != nil {
		fmt.Println("ОШИБКА при создании:", err)
	} else {
		fmt.Println("УСПЕХ! Задача создана:")
		fmt.Println("  ID пользователя:", createdTask.UserID)
		fmt.Println("  Заголовок:", createdTask.Title)
		fmt.Println("  Статус по умолчанию:", createdTask.Status)
		fmt.Println("  Время создания:", createdTask.CreatedAt.Format("15:04:05 02.01.2006"))
	}

	// -------------------------------------------------------------
	// ТЕСТ 2: Проверка защиты от дурака (пустой заголовок)
	// -------------------------------------------------------------
	fmt.Println("\n[ТЕСТ 2] Проверяем задачу с пустым заголовком...")
	badTask := domain.Task{
		UserID: 1,
		Title:  "",
	}

	_, err = taskService.CreateTask(badTask)
	if err != nil {
		fmt.Println("ОТЛИЧНО! Сервис перехватил ошибку:", err)
	} else {
		fmt.Println("БАГ! Сервис почему-то пропустил пустую задачу!")
	}

	//тест для метода с выдачей
	fmt.Println("\n[ТЕСТ 3] Проверяем GetTaskID...")
	testTask := domain.Task{
		ID:     2,
		UserID: 2,
		Title:  "тест поиска",
	}

	foundID, err := taskService.GetTaskID(testTask.ID)
	if err != nil {
		fmt.Println(err, "ошибка")
	} else {
		fmt.Println("успешно, метод вернул ID", foundID)
	}

	fmt.Println("\n[ТЕСТ 4] Првоеряем UpdateTask...")
	testTask = domain.Task{
		ID: 1,
		UserID: 2,
		Title: "тест update",
		Description: "sixseven",
		Status: domain.TaskStatusCancelled,
	}
	testTask1, err := taskService.UpdateStatus(1,2,domain.TaskStatusCancelled)
	if err != nil{
		fmt.Println("error", err)
	}else{
		fmt.Println("Успешно, статус задачи стал: ", testTask1.Status)
	}

	//проверяем новый метод
	fmt.Println("\n[ТЕСТ 5] Проверяем DeleteTask...")
	testTask2 := domain.Task{
		ID: 3,
		Title: "test delete task",
		UserID: 3,
		Description: "test",
	}
	err = taskService.DeleteTask(testTask2.ID, testTask2.UserID)
	if err != nil{
		fmt.Println("error", err)
	}else{
		fmt.Println("Успех, задача удалена")
	}

	//проверка последнего метода
	fmt.Println("\n[Тест 6] Првоеряем GetList...")
	testTask3 := domain.Task{
		ID: 4,
		Title: "test get list",
		UserID: 4,
		Description: "test",
	}
	testTaskList, err := taskService.GetList(testTask3.UserID)
	if err != nil{
		fmt.Println("error", err)
	}else{
		fmt.Println("Успех, список выдан", testTaskList)
	}

	fmt.Println("\n[Расписание]")

	//тестим метод AddSlot
	fmt.Println("\n[Тест 1] Проверяем AddSlot...")
	scheduleTest1 := domain.ScheduleSlot{
		UserID: 1,
		Subject: "МДК 05.02",
		Room: "1000-7",
		StartTime: "09:00",
		EndTime: "10:30",
		DayOfWeek: 3,
		Type: domain.SlotTypeLecture,
		Parity: domain.WeekParityBoth,
	}
	testSchedule1, err := scheduleService.AddSlot(scheduleTest1)
	if err != nil{
		fmt.Println("error: ", err)
	}else{
		fmt.Println("Успех!",
	"\n id пользователя: ", scheduleTest1.UserID,
	"\n название предмета: ", scheduleTest1.Subject,
	"\n кабинет: ", scheduleTest1.Room,
	"\n день недели: ", scheduleTest1.DayOfWeek,
	"\n тип: ", scheduleTest1.Type,
	"\n четность: ", scheduleTest1.Parity,
	"\n время начала: ", scheduleTest1.StartTime,
	"\n время окончания: ", scheduleTest1.EndTime,)
	fmt.Println(testSchedule1)
	}
}
