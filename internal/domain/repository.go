package domain

type TaskRepository interface{
	CreateTask()
	GetTaskID()
	GetList()
	UpdateTask()
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