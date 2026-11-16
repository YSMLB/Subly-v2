package domain

import (
	"errors"
	"time"
)

var (
	ErrRateLimited       = errors.New("слишком частые запросы, подождите минуту")
	ErrInvalidCode       = errors.New("неверный или истекший код")
	ErrUserAlreadyExists = errors.New("пользователь с таким email уже существует")
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	GroupName    string    `json:"group_name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type VerificationCode struct {
	ID        int
	Email     string
	Code      string
	ExpiresAt time.Time
	Attempts  int
}

type EmailSender interface {
	SendVerificationCode(toEmail, code string) error
}

type AuthRepository interface {
	SaveVerificationCode(email, code string, ttl time.Duration) error
	GetValidCode(email string) (*VerificationCode, error)
	IncrementAttempts(id int) error
	DeleteCode(email string) error

	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	IsCodeRequestedRecently(email string) (bool, error)
}
