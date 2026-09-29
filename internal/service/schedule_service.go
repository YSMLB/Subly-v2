package service

import (
	"errors"
	"fmt"
	"subly-v2/internal/domain"
)

type ScheduleService struct{
	Repo domain.ScheduleRepository
}

//type ScheduleSlot struct{
//	ID int `json:"id"`
//	UserID int `json:"userID"`
//	Subject string `json:"subject"`
//	Teacher string `json:"teacher"`
//	Room string `json:"room"`
//	DayOfWeek int `json:"dayOfWeek"`
//	StartTime string `json:"startTime"`
//	EndTime string `json:"endTime"`
//	Type SlotType `json:"slotType"`
//	Parity WeekParity `json:"weekParity"`
//	IsCancelled bool `json:"isCancelled"`
//}

func (as *ScheduleService)AddSlot(schedule domain.ScheduleSlot) (*domain.ScheduleSlot, error){
	if schedule.UserID <= 0{
		return nil, errors.New("некорректные данные студента")
	}
	if schedule.Subject == ""{
		return nil, errors.New("название предмета не может быть пустым")
	}
	if schedule.StartTime == "" || schedule.EndTime == ""{
		return nil, errors.New("Время начала или конца занятия должно быть заполнено")
	}
	if schedule.DayOfWeek > 7 || schedule.DayOfWeek <=0 {
		return nil, errors.New("День недели должен быть корректным")
	}
	if schedule.IsCancelled == true{
		return nil, errors.New("поле не может быть отмененным")
	}

	err := as.Repo.AddSlot(schedule)
	if err != nil{
		fmt.Println("error: ", err)
		return nil, err
	}

	return &schedule, nil
}

func (gst *ScheduleService) GetScheduleToday() error{
	return nil
}