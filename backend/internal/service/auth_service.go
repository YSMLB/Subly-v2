package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"subly-v2/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   domain.AuthRepository
	mailer domain.EmailSender
	secret string
}

func NewAuthService(repo domain.AuthRepository, mailer domain.EmailSender, secret string) *AuthService {
	return &AuthService{
		repo:   repo,
		mailer: mailer,
		secret: secret,
	}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Code      string `json:"code"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	GroupName string `json:"group_name"`
}

func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func (a *AuthService) RequestCode(email string) error {
	isRecent, err := a.repo.IsCodeRequestedRecently(email)
	if err != nil {
		return err
	}
	if isRecent {
		return domain.ErrRateLimited
	}

	code, err := GenerateOTP()
	if err != nil {
		return err
	}

	err = a.repo.SaveVerificationCode(email, code, 10*time.Minute)
	if err != nil {
		return err
	}

	if err := a.mailer.SendVerificationCode(email, code); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Register(req RegisterRequest) (string, *domain.User, error) {
	existing, err := a.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", nil, err
	}
	if existing != nil {
		return "", nil, domain.ErrUserAlreadyExists
	}

	validCode, err := a.repo.GetValidCode(req.Email)
	if err != nil {
		return "", nil, err
	}
	if validCode == nil {
		return "", nil, errors.New("код подтверждения истек или не существует")
	}

	if validCode.Code != req.Code {
		_ = a.repo.IncrementAttempts(validCode.ID)
		if validCode.Attempts+1 >= 5 {
			_ = a.repo.DeleteCode(req.Email)
		}
		return "", nil, domain.ErrInvalidCode
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка хэширования пароля: %w", err)
	}

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		GroupName:    req.GroupName,
		Role:         "student",
		CreatedAt:    time.Now(),
	}

	if err := a.repo.CreateUser(user); err != nil {
		return "", nil, err
	}

	_ = a.repo.DeleteCode(req.Email)

	token, err := a.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return token, user, nil
}

func (a *AuthService) Login(email, password string) (string, *domain.User, error) {
	user, err := a.repo.GetUserByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("неверный email или пароль")
	}

	token, err := a.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (a *AuthService) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}
