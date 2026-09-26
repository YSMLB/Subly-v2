package domain

type SlotType string

const(
	SlotTypeLecture SlotType = "lecture"
	SlotTypePractice SlotType = "practice"
	SlotTypeLab SlotType= "lab"
)

type WeekParity string

const(
	WeekParityBoth WeekParity = "both"
	WeekParityOdd WeekParity = "odd"
	WeekParityEven WeekParity = "even"
)

type ScheduleSlot struct{
	ID int `json:"id"`
	UserID int `json:"userID"`
	Subject string `json:"subject"`
	Teacher string `json:"teacher"`
	Room string `json:"room"`
	DayOfWeek int `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Type SlotType `json:"slotType"`
	Parity WeekParity `json:"weekParity"`
	IsCancelled bool `json:"isCancelled"`
}