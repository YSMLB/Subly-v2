package domain

import "time"

type ChangeStatus string

const(
	ChangeStatusPending ChangeStatus = "pending"
	ChangeStatusApproved ChangeStatus = "approved"
	ChangeStatusRejected ChangeStatus = "rejected"
)

type ChangeAction string

const(
	ChangeActionCancel ChangeAction = "cancel"
	ChangeActionRescheduleTime ChangeAction = "reschedule_time"
	ChangeActionChangeRoom ChangeAction = "change_room"
	ChangeActionAddSlot ChangeAction = "add_slot"
)

type ProposedChange struct{
	ID int `json:"id"`
	UserID int `json:"userID"`
	SlotID int `json:"slotID"`
	Action ChangeAction `json:"changeAction"`
	Status ChangeStatus `json:"changeStatus"`
	RawText string `json:"rawText"`
	NewStartTime string `json:"newStartTime"`
	NewEndTime string `json:"newEndTime"`
	NewRoom string `json:"newRoom"`
	CreatedAdd time.Time `json:"createdAdd"`
}