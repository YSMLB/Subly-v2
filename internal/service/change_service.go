package service

import (
	"errors"
	"subly-v2/internal/domain"
	"subly-v2/internal/ingestor"
	"time"
)

type ChangeService struct {
	Repo domain.ChangeRepository
	RepoSchedule domain.ScheduleRepository
}

func (sc *ChangeService) SaveChanges(changes domain.ProposedChange) (*domain.ProposedChange, error) {
	changes.Status = domain.ChangeStatusPending

	if changes.UserID <= 0 {
		return nil, errors.New("Некорректный UserID")
	}
	if changes.Action != domain.ChangeActionAddSlot && changes.Action != domain.ChangeActionCancel && changes.Action != domain.ChangeActionChangeRoom && changes.Action != domain.ChangeActionRescheduleTime {
		return nil, errors.New("Не зафиксирован допустимый доменный тип планируемого действия")
	}

	if changes.Status != domain.ChangeStatusPending {
		return nil, errors.New("Некорректный статус")
	}

	if changes.RawText == "" {
		return nil, errors.New("Некорректно передано изменение")
	}
	if changes.NewEndTime == "" && changes.NewStartTime == "" && changes.NewRoom == "" {
		return nil, errors.New("Нет определенного действия")
	}

	changes.CreatedAdd = time.Now()
	changesSave, err := sc.Repo.SaveChanges(changes)
	if err != nil {
		return nil, err
	}

	return changesSave, nil
}

func (gpbuid *ChangeService) GetPendingByUserID(userID int) ([]domain.ProposedChange, error) {
	if userID <= 0 {
		return nil, errors.New("Некорректный UserID")
	}
	waitChange, err := gpbuid.Repo.GetPendingByUserID(userID)
	if err != nil {
		return nil, err
	}

	return waitChange, nil
}

func (gcbid *ChangeService)GetChangeByID(id int) (*domain.ProposedChange, error){
	if id <= 0{
		return nil, errors.New("Некорректный id")
	}

	changes, err := gcbid.Repo.GetChangeByID(id)

	if err != nil{
		return nil, err
	}
	return changes, nil
}

func (usc *ChangeService) UpdateStatusChanges(slotID int, userID int, newStatus domain.ChangeStatus) (*domain.ProposedChange, error){
	if userID <= 0 {
		return nil, errors.New("Некорректный userID")
	} else if slotID <= 0 {
		return nil, errors.New("Некорректный slotID")
	}

	if newStatus == domain.ChangeStatusPending{
		return nil, errors.New("Статус должен быть: Принято или Отклонено")
	}

	changes, err := usc.Repo.GetChangeByID(slotID)
	if err != nil{
		return nil, err
	}

	if changes.UserID != userID{
		return nil, errors.New("Вы не можете отклонить чужое изменение")
	}

	if changes.Status != domain.ChangeStatusPending{
		return nil, errors.New("Это изменение не активно")
	}

	if newStatus == domain.ChangeStatusApproved{
		slot, err := usc.RepoSchedule.GetSlotByID(changes.SlotID)
		if err != nil{
			return nil, err
		}

		if slot == nil{
			return nil, errors.New("Слот не найден")
		}

		switch changes.Action{
		case domain.ChangeActionCancel:
			slot.IsCancelled = true

		case domain.ChangeActionChangeRoom:
			slot.Room = changes.NewRoom

		case domain.ChangeActionRescheduleTime:
			slot.StartTime = changes.NewStartTime
			slot.EndTime = changes.NewEndTime
		}

		err = usc.RepoSchedule.UpdateTaskSchedule(*slot)
		if err != nil{
			return nil, errors.New("UpdateTaskSchedule вернул ошибку")
		}

	}


	changes.Status = newStatus

	finallyChanges, err := usc.Repo.UpdateStatusChanges(slotID, userID, newStatus)

	return finallyChanges, nil
}


func (s *ChangeService)ImportReplacementsFromText(userID int, group string, rawText string) ([]domain.ProposedChange, error){
	if userID <= 0{
		return nil, errors.New("userID должен иметь значение отличное от нуля")
	}

	if group == ""{
		return nil, errors.New("поле группы не должно юыть пустым")
	}

	if rawText == ""{
		return nil, errors.New("поле с изменением не должно быть пустым")
	}

	sliceChangeParse , err := ingestor.MaxChangeParser(rawText, group)
	if err != nil{
		return nil, errors.New("MaxChangeParser вернул ошибку")
	}

	scheduleSlotToday, err := s.RepoSchedule.GetScheduleToday(userID, sliceChangeParse.DayOfWeek)
	if err != nil{
		return nil, errors.New("GetScheduleToday вернул ошибку")
	}

	result := make([]domain.ProposedChange, 0)

	for i := 0; i < len(sliceChangeParse.ListSlots); i++{
		res := &domain.ScheduleSlot{}
		change := &domain.ProposedChange{}
		n := &sliceChangeParse.ListSlots[i]

		for _ ,v := range scheduleSlotToday{
			if v.Subject == n.Subject{
				res.ID = v.ID
				break
			}
		}

		if n.IsCancelled == true{
			change.Action = domain.ChangeActionCancel
		}else if res.ID > 0{
			change.Action = domain.ChangeActionChangeRoom
		} else{
			change.Action = domain.ChangeActionAddSlot
		}

		change = &domain.ProposedChange{
			UserID: userID,
			SlotID: res.ID,
			Action: change.Action,
			Status: domain.ChangeStatusPending,
			NewRoom: n.Room,
			RawText: n.Subject,
		}

		preres, err := s.Repo.SaveChanges(*change)
		if err != nil{
			return nil, err
		}

		result = append(result, *preres)
	}


	return result, nil
}
