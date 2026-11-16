package service

import (
	"crypto/rand"
	"math/big"
	"subly-v2/internal/domain"
	"time"
	"fmt"
)

type AuthService struct {
    repo   domain.AuthRepository
    mailer domain.EmailSender
    secret string
}

func GenerateOTP() (string, error){
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil{
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()+100000), nil//формула, которую мне дал гпт
}

func (a *AuthService)RequestCode(email string) error{
	//проверяем действителен ли код(меньше 60 сек назад или нет)
	isRecent, err := a.repo.IsCodeRequestedRecently(email)
	if err != nil{
		return err
	}
	if isRecent{
		return domain.ErrRateLimited
	}

	//генерим кол
	code, err := GenerateOTP()
	if err != nil{
		return err
	}

	//сохраняем на 10 мин
	err = a.repo.SaveVerificationCode(email, code, 10 * time.Minute)
	if err != nil{
		return err
	}

	//отправка
	if err := a.mailer.SendVerificationCode(email, code); err != nil{
		return err
	}

	return nil
}