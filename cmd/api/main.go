package main

import (
	"fmt"
	"subly-v2/internal/domain"
	"subly-v2/internal/service"
	"time"
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
	fmt.Println("--- Запуск проверки логики Subly ---")

	// 1. Создаем заглушку репозитория
	mockRepo := &MockTaskRepo{}

	// 1. Создали заглушку репы
	mockRepoSchedule := &MockScheduleRepo{}

	// 1. Создали заглушку репы
	mockRepoChange := &MockChangeRepo{}

	// 2. Инициализируем наш сервис задач
	taskService := &service.TaskService{
		Repo: mockRepo,
	}

	scheduleService := &service.ScheduleService{
		Repo: mockRepoSchedule,
	}

	changeService := &service.ChangeService{
		Repo: mockRepoChange,
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


	fmt.Println("\n[Изменение]")//change
	
	fmt.Println("\n [Тест 1] Проверяем SaveChange...")//savechange
	changeTest1 := domain.ProposedChange{
				ID: 2,
		UserID: 3,
		SlotID: 4,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "МДК 05.02",
		NewStartTime: "12:10",
		NewEndTime: "13:40",
		NewRoom: "401",
		CreatedAdd: time.Now(),
	}

	testChange1, err := changeService.SaveChanges(changeTest1)
	if err != nil{
		fmt.Println("error: ", err)
	}else{
		fmt.Println(
			"\nID: ", changeTest1.ID,
			"\nUserID: ", changeTest1.UserID,
			"\nSlotID: ", changeTest1.SlotID,
			"\nAction: ", changeTest1.Action,
			"\nStatus: ", changeTest1.Status,
			"\nRawText: ", changeTest1.RawText,
			"\nNewStartTime: ", changeTest1.NewStartTime,
			"\nNewEndTime: ", changeTest1.NewEndTime,
			"\nNewRoom: ", changeTest1.NewRoom,
			"\nCreatedAdd: ", changeTest1.CreatedAdd,
		)
		fmt.Println(testChange1)
	}

	fmt.Println("[Тест 2] Проверяем GetPendingByUserID...")
	changeTest2 := domain.ProposedChange{
		ID: 3,
		UserID: 4,
		SlotID: 5,
		Action: domain.ChangeActionAddSlot,
		Status: domain.ChangeStatusPending,
		RawText: "МДК 05.02",
		NewStartTime: "12:10",
		NewEndTime: "13:40",
		NewRoom: "401",
		CreatedAdd: time.Now(),
	}

	testChange2, err := changeService.GetPendingByUserID(changeTest2.UserID,)
	if err != nil{
		fmt.Println("error: ", err)
	}else{
		fmt.Println(
			"\nID: ", changeTest2.ID,
			"\nUserID: ", changeTest2.UserID,
			"\nSlotID: ", changeTest2.SlotID,
			"\nAction: ", changeTest2.Action,
			"\nStatus: ", changeTest2.Status,
			"\nRawText: ", changeTest2.RawText,
			"\nNewStartTime: ", changeTest2.NewStartTime,
			"\nNewEndTime: ", changeTest2.NewEndTime,
			"\nNewRoom: ", changeTest2.NewRoom,
			"\nCreatedAdd: ", changeTest2.CreatedAdd,
		)
		fmt.Println(testChange2)
	}

		fmt.Println("[Тест 3] Проверяем GetChangeByID...")

		changeTest3 := domain.ProposedChange{
			ID: 5,
			UserID: 4,
			SlotID: 2,
			Action: domain.ChangeActionAddSlot,
			Status: domain.ChangeStatusPending,
			RawText: "geography",
			NewStartTime: "15:40",
			NewEndTime: "16:50",
			NewRoom: "409",
			CreatedAdd: time.Now(),
		}
		testChange3, err := changeService.GetChangeByID(changeTest3.ID)
		if err != nil{
			fmt.Println("error: ", err)
		}else{
			fmt.Println(
			"\nID: ", changeTest3.ID,
			"\nUserID: ", changeTest3.UserID,
			"\nSlotID: ", changeTest3.SlotID,
			"\nAction: ", changeTest3.Action,
			"\nStatus: ", changeTest3.Status,
			"\nRawText: ", changeTest3.RawText,
			"\nNewStartTime: ", changeTest3.NewStartTime,
			"\nNewEndTime: ", changeTest3.NewEndTime,
			"\nNewRoom: ", changeTest3.NewRoom,
			"\nCreatedAdd: ", changeTest3.CreatedAdd,
		)
		fmt.Println(testChange3)
		}

		fmt.Println("[Тест 3] Проверяем UpdateStatusChanges...")

		changeTest4 := domain.ProposedChange{
			ID: 4,
			UserID: 4,
			SlotID: 6,
			Action: domain.ChangeActionAddSlot,
			Status: domain.ChangeStatusPending,
			RawText: "Высшая математика",
			NewStartTime: "13:00",
			NewEndTime: "13:30",
			NewRoom: "33",
			CreatedAdd: time.Now(),
		}
		newStatus := domain.ChangeStatusApproved

		testChange4, err := changeService.UpdateStatusChanges(changeTest4.SlotID, changeTest4.UserID, newStatus)

		if err != nil{
			fmt.Println("error: ", err)
		}else{
			fmt.Println(
			"\nID: ", changeTest4.ID,
			"\nUserID: ", changeTest4.UserID,
			"\nSlotID: ", changeTest4.SlotID,
			"\nAction: ", changeTest4.Action,
			"\nStatus: ", changeTest4.Status,
			"\nRawText: ", changeTest4.RawText,
			"\nNewStartTime: ", changeTest4.NewStartTime,
			"\nNewEndTime: ", changeTest4.NewEndTime,
			"\nNewRoom: ", changeTest4.NewRoom,
			"\nCreatedAdd: ", changeTest4.CreatedAdd,
		)
		fmt.Println(testChange4)
		}
	fmt.Println("===Конец тестов===")
}
