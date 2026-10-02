package service

import (
	"errors"
	"subly-v2/internal/domain"
	"time"
)

type ChangeService struct {
	Repo domain.ChangeRepository
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
