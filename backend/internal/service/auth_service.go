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

// GenerateOTP генерирует криптостойкий 6-значный код (от 100000 до 999999)
func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// RequestCode отправляет код подтверждения на email с проверкой rate limit
func (a *AuthService) RequestCode(email string) error {
	// 1. Проверяем действителен ли код (меньше 60 сек назад или нет)
	isRecent, err := a.repo.IsCodeRequestedRecently(email)
	if err != nil {
		return err
	}
	if isRecent {
		return domain.ErrRateLimited
	}

	// 2. Генерируем 6-значный код
	code, err := GenerateOTP()
	if err != nil {
		return err
	}

	// 3. Сохраняем код в БД на 10 минут
	err = a.repo.SaveVerificationCode(email, code, 10*time.Minute)
	if err != nil {
		return err
	}

	// 4. Отправляем письмо через mailer
	if err := a.mailer.SendVerificationCode(email, code); err != nil {
		return err
	}

	return nil
}

// Register валидирует OTP-код, хэширует пароль, сохраняет пользователя и выдаёт JWT
func (a *AuthService) Register(req RegisterRequest) (string, *domain.User, error) {
	// 1. Проверяем, существует ли уже пользователь с таким email
	existing, err := a.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", nil, err
	}
	if existing != nil {
		return "", nil, domain.ErrUserAlreadyExists
	}

	// 2. Ищем действующий код в БД
	validCode, err := a.repo.GetValidCode(req.Email)
	if err != nil {
		return "", nil, err
	}
	if validCode == nil {
		return "", nil, errors.New("код подтверждения истек или не существует")
	}

	// 3. Сверяем введённый код
	if validCode.Code != req.Code {
		_ = a.repo.IncrementAttempts(validCode.ID)
		if validCode.Attempts+1 >= 5 {
			_ = a.repo.DeleteCode(req.Email)
		}
		return "", nil, domain.ErrInvalidCode
	}

	// 4. Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка хэширования пароля: %w", err)
	}

	// 5. Создаём сущность пользователя
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

	// 6. Удаляем использованный код подтверждения
	_ = a.repo.DeleteCode(req.Email)

	// 7. Генерируем JWT токен на 30 дней
	token, err := a.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return token, user, nil
}

// generateToken выпускает JWT токен для пользователя
func (a *AuthService) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}
