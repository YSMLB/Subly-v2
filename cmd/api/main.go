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

	fmt.Println("\n[Расписание]")//schedule

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


	//test method Get...
	fmt.Println("\n[Тест 2] Проверяем GetScheduleToday...")

	scheduleTest2 := domain.ScheduleSlot{
		UserID: 2,
		Subject: "МДК 05.01",
		Room: "402",
		StartTime: "08:30",
		EndTime: "10:00",
		DayOfWeek: 5,
		Type: domain.SlotTypeLecture,
		Parity: domain.WeekParityBoth,
	}

	testSchedule2, err := scheduleService.GetScheduleToday(scheduleTest2.UserID, scheduleTest2.DayOfWeek, scheduleTest2.Parity)
	if err != nil{
		fmt.Println("error: ", err)
	}else{
		fmt.Println("Успех!",
	"\n id пользователя: ", testSchedule2.UserID,
	"\n название предмета: ", testSchedule2.Subject,
	"\n кабинет: ", testSchedule2.Room,
	"\n день недели: ", testSchedule2.DayOfWeek,
	"\n тип: ", testSchedule2.Type,
	"\n четность: ", testSchedule2.Parity,
	"\n время начала: ", testSchedule2.StartTime,
	"\n время окончания: ", testSchedule2.EndTime,)
	fmt.Println(testSchedule2)
	}

	fmt.Println("\n[Тест 3] Проверяем UpdateTaskSchedule...")
	scheduleTest3 := domain.ScheduleSlot{
		ID: 1,
		UserID: 3,
		Subject: "ОАиП",
		Room: "407",
		StartTime: "08:30",
		EndTime: "10:00",
		DayOfWeek: 6,
		Type: domain.SlotTypeLecture,
		Parity: domain.WeekParityBoth,
	}

	err = scheduleService.UpdateTaskSchedule(scheduleTest3)
	if err != nil{
		fmt.Println("ошибочка: ", err)
	}else{
		fmt.Println("Успех!",
	"\n id пользователя: ", scheduleTest3.UserID,
	"\n название предмета: ", scheduleTest3.Subject,
	"\n кабинет: ", scheduleTest3.Room,
	"\n день недели: ", scheduleTest3.DayOfWeek,
	"\n тип: ", scheduleTest3.Type,
	"\n четность: ", scheduleTest3.Parity,
	"\n время начала: ", scheduleTest3.StartTime,
	"\n время окончания: ", scheduleTest3.EndTime,)
	}
}
