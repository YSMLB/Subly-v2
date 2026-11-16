package service

import "subly-v2/internal/domain"

type AuthService struct {
    repo   domain.AuthRepository
    mailer domain.EmailSender
    secret string
}

func (a *AuthService)RequestCode(email string) error{
	
}