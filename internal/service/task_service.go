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