package service

import (
	"errors"
	"subly-v2/internal/domain"
	"time"
)

type ChangeService struct{
	Repo domain.ChangeRepository
}

func (sc *ChangeService)SaveChanges(changes domain.ProposedChange) (*domain.ProposedChange, error){
	changes.Status = domain.ChangeStatusPending
	if changes.ID <= 0{
		return nil, errors.New("Некорректный id")
	}
	if changes.SlotID<= 0{
		return nil, errors.New("Некорректный SlotID")
	}
	if changes.UserID <= 0{
		return nil, errors.New("Некорректный UserID")
	}
	if changes.Action != domain.ChangeActionAddSlot && changes.Action != domain.ChangeActionCancel && changes.Action != domain.ChangeActionChangeRoom && changes.Action != domain.ChangeActionRescheduleTime{
		return nil, errors.New("Не зафиксирован допустимый доменный тип планируемого действия")
	}

	if changes.Status != domain.ChangeStatusPending{
		return nil, errors.New("Некорректный статус")
	}

	if changes.RawText == ""{
		return nil, errors.New("Некорректно передано изменение")
	}
	if changes.NewEndTime == "" || changes.NewStartTime == "" || changes.NewRoom == ""{
		return nil, errors.New("Нет определенного действия")
	}

	changesSave, err := sc.Repo.SaveChanges(changes)
	if err != nil{
		return nil, err
	}
	changesSave.CreatedAdd = time.Now()

	return changesSave, nil
}