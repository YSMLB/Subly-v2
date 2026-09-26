package domain

type TaskRepository interface{
	CreateTask(Task) error
	GetTaskID(id int) (*Task, error)
	GetList()
	UpdateTask(task *Task) error
	DeleteTask()
}

type ScheduleRepository interface{
	GetScheduleToday()
	UpdateTaskSchedule()
}

type ChangeRepository interface{
	SaveChanges()
	GetPendingByUserID()
	UpdateStatusChanges()
}