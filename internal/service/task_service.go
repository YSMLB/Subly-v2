package service

import (
	"errors"
	"subly-v2/internal/domain"
	"time"
)

type TaskService struct{
	Repo domain.TaskRepository
}

func (s *TaskService) CreateTask(task domain.Task) (*domain.Task, error){
	if task.Title == ""{
		return nil, errors.New("заголовок задачи не может быть пустым")
	}

	if task.DueDate.Before(time.Now()) && !task.DueDate.IsZero(){
		return nil, errors.New("дедлайн не может быть в прошлом")
	}

	if task.Status == ""{
		task.Status = domain.TaskStatusPending
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	err := s.Repo.CreateTask(task)
	if err != nil{
		return nil, err
	}

	return &task,  nil
}

func (us *TaskService) UpdateStatus(id int, userID int, newStatus domain.TaskStatus) (*domain.Task, error) {
	if id <= 0 {
		return nil, errors.New("Такого пользователя не существует!")
	}
	if userID <= 0{
		return nil, errors.New("Такой задачи не существует!")
	}

	task, err := us.Repo.GetTaskID(id)
	if err != nil{
		return nil, err
	}
	

	if task.UserID != userID{
		return nil, errors.New("нет прав на изменение чужой задачи!")
	}
	if task.Status == domain.TaskStatusCancelled && newStatus == domain.TaskStatusCompleted{
		return nil, errors.New("Нельзя завершить отмененную задачу!")
	}

	task.Status = newStatus
	err = us.Repo.UpdateTask(task)
	if err != nil{
		return nil, err
	}

	return task, nil 
}

func (gtid *TaskService) GetTaskID(id int) (*domain.Task ,error){
	if id == 0{
		return nil, errors.New("Такого пользователя не существует!")
	}
	
	task, err := gtid.Repo.GetTaskID(id)
	if err != nil{
		return nil, err
	}
	return task, nil
	//if task.UserID == 0{
	//	return 0, errors.New("Такого пользователя не существует!")
	//}
//
	//err := gtid.Repo.GetTaskID(task)
	//if err != nil{
	//	return 0, err
	//}
	//return task.ID, nil
}

func (dt *TaskService) DeleteTask(taskID int, userID int) error{
	if taskID <= 0{
		return errors.New("Такой задачи не существует!")
	}
	if userID <= 0 {
		return errors.New("Такого пользователя не существует!")
	}
	
	task, err := dt.Repo.GetTaskID(taskID)
	if err != nil{
		return err
	}
	if task.UserID != userID{
		return errors.New("Вы не можете удалить чужую задачу!")
	}

	err = dt.Repo.DeleteTask(task.ID)
	if err != nil{
		return err
	}

	return nil
}