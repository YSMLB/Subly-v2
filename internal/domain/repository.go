package domain



type TaskRepository interface {
	CreateTask(Task) error
	GetTaskID(id int) (*Task, error)
	GetList(userID int) ([]Task, error)
	UpdateTask(task *Task) error
	DeleteTask(taskID int) error
}

type ScheduleRepository interface {
	AddSlot(schedule ScheduleSlot) error
	GetScheduleToday(userID int, dayOfWeek int) (*ScheduleSlot, error)
	UpdateTaskSchedule(schedule ScheduleSlot) error
}

type ChangeRepository interface {
	SaveChanges(changes ProposedChange) (*ProposedChange, error)
	GetPendingByUserID(userID int) ([]ProposedChange, error)
	UpdateStatusChanges()
}