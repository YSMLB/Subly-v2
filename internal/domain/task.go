package domain

import "time"

//type TaskStatus struct{
//	Pending string`json:"pending"`
//	InProgress string
//}
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskPriority string

const (
	TaskPriorityLow TaskPriority="Low"
	TaskPriorityMedium TaskPriority="Medium"
	TaskPriorityHigh TaskPriority="High"
)

type Task struct{
	ID int `json:"id"`
	UserID int `json:"userID"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status TaskStatus `json:"status"`
	Priority TaskPriority `json:"priority"`
	DueDate time.Time `json:"dueDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
}