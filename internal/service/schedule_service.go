package service

import (
	"errors"
	"fmt"
	"subly-v2/internal/domain"
	"subly-v2/internal/ingestor"
)

type ScheduleService struct {
	Repo domain.ScheduleRepository
}

func (as *ScheduleService) AddSlot(schedule domain.ScheduleSlot) (*domain.ScheduleSlot, error) {
	if schedule.UserID <= 0 {
		return nil, errors.New("некорректные данные студента")
	}
	if schedule.Subject == "" {
		return nil, errors.New("название предмета не может быть пустым")
	}
	if schedule.StartTime == "" || schedule.EndTime == "" {
		return nil, errors.New("Время начала или конца занятия должно быть заполнено")
	}
	if schedule.DayOfWeek > 7 || schedule.DayOfWeek <= 0 {
		return nil, errors.New("День недели должен быть корректным")
	}
	if schedule.IsCancelled == true {
		return nil, errors.New("поле не может быть отмененным")
	}

	err := as.Repo.AddSlot(schedule)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	return &schedule, nil
}

func (gst *ScheduleService) GetScheduleToday(userID int, dayOfWeek int, parity domain.WeekParity) ([]domain.ScheduleSlot, error) {
	if userID <= 0 {
		return nil, errors.New("Некорректный userID")
	}
	if dayOfWeek <= 0 || dayOfWeek > 7 {
		return nil, errors.New("Некорректный день недели")
	}
	if parity != domain.WeekParityBoth && parity != domain.WeekParityEven && parity != domain.WeekParityOdd {
		return nil, errors.New("Некорректная четность недели")
	}

	schedule, err := gst.Repo.GetScheduleToday(userID, dayOfWeek)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	filtered := make([]domain.ScheduleSlot, 0)
	for _, slot := range schedule {
		if slot.Parity != domain.WeekParityBoth && parity != domain.WeekParityBoth && slot.Parity != parity {
			continue
		}
		if slot.Subject == "" {
			continue
		}
		filtered = append(filtered, slot)
	}

	return filtered, nil
}

func (uts *ScheduleService) UpdateTaskSchedule(schedule domain.ScheduleSlot) error {
	if schedule.UserID <= 0 {
		return errors.New("Некорректный UserID")
	}

	if schedule.Subject == "" {
		return errors.New("Некорректное название")
	}

	if schedule.StartTime == "" || schedule.EndTime == "" {
		return errors.New("Занятие должно иметь временные рамки")
	}

	if schedule.DayOfWeek <= 0 || schedule.DayOfWeek > 6 {
		return errors.New("День должен быть корректный")
	}

	if schedule.ID <= 0 {
		return errors.New("Такой пары не существует")
	}

	err := uts.Repo.UpdateTaskSchedule(schedule)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScheduleService) ImportScheduleFromFile(userID int, filePath string) error {
	parseData, err := ingestor.ParseSchedule(filePath)
	if err != nil {
		return err
	}

	for _, i := range parseData {
		i.UserID = userID

		err := s.Repo.AddSlot(i)
		if err != nil {
			return err
		}
	}

	return nil
}
