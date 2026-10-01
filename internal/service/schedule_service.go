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

func (gst *ScheduleService) GetScheduleToday(userID int, dayOfWeek int, parity domain.WeekParity) (*domain.ScheduleSlot, error){
	if userID <= 0{
		return nil, errors.New("Некорректный userID")
	}
	if dayOfWeek <= 0 || dayOfWeek > 7{
		return nil, errors.New("Некорректный день недели")
	}
	if parity != domain.WeekParityBoth && parity != domain.WeekParityEven && parity != domain.WeekParityOdd{
		return nil, errors.New("Некорректная четность недели")
	}

	schedule, err := gst.Repo.GetScheduleToday(userID, dayOfWeek)

	if err != nil{
		fmt.Println("error: ", err)
		return nil, err
	}

	if schedule.Parity != domain.WeekParityBoth || schedule.Parity != parity {
		return nil, errors.New("На этой неделе нет таких пар")
	}

	if schedule.Subject == ""{
		return nil, nil
	}
	

	return schedule, nil
}

func (uts *ScheduleService) UpdateTaskSchedule(schedule domain.ScheduleSlot) error{
	if schedule.UserID <= 0{
		return errors.New("Некорректный UserID")
	}

	if schedule.Subject == ""{
		return errors.New("Некорректное название")
	}

	if schedule.StartTime == "" || schedule.EndTime == ""{
		return errors.New("Занятие должно иметь временные рамки")
	}

	if schedule.DayOfWeek <= 0 || schedule.DayOfWeek > 6{
		return errors.New("День должен быть корректный")
	}
	
	if schedule.ID <= 0{
		return errors.New("Такой пары не существует")
	}

	err := uts.Repo.UpdateTaskSchedule(schedule)

	if err != nil{
		return err
	}
	return nil
}







//Even (и́вен) — Чётная неделя (в вузах часто «знаменатель»)
//💡 Лайфхак: в слове even ровно 4 буквы (4 — чётное число) 
//→
//→ Чётная!
//
//2. Odd (одд) — Нечётная неделя (в вузах «числитель»)
//💡 Лайфхак: в слове odd ровно 3 буквы (3 — нечётное число) 
//→
//→ Нечётная!
//
//3. Both (бо́ус) — Обе недели (и та, и другая)
//Означает, что занятие проходит каждую неделю подряд, без разницы — числитель сейчас или знаменатель.