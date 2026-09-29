package domain

type TaskRepository interface{
	CreateTask(Task) error
	GetTaskID(id int) (*Task, error)
	GetList(userID int) ([]Task, error)
	UpdateTask(task *Task) error
	DeleteTask(taskID int) error
}

type ScheduleRepository interface{
	AddSlot(schedule ScheduleSlot)  error
	GetScheduleToday()
	UpdateTaskSchedule()
}

type ChangeRepository interface{
	SaveChanges()
	GetPendingByUserID()
	UpdateStatusChanges()
}