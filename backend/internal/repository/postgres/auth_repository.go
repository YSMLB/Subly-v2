package postgres

import (
	"context"
	"errors"
	"time"

	"subly-v2/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) SaveVerificationCode(email, code string, ttl time.Duration) error {
	ctx := context.Background()

	deleteQuery := "DELETE FROM verification_codes WHERE email = $1"
	_, err := r.pool.Exec(ctx, deleteQuery, email)
	if err != nil {
		return err
	}

	// Вставляем новый код с временем жизни
	expiresAt := time.Now().Add(ttl)
	insertQuery := "INSERT INTO verification_codes (email, code, expires_at, attempts, created_at) VALUES ($1, $2, $3, 0, NOW())"
	_, err = r.pool.Exec(ctx, insertQuery, email, code, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetValidCode(email string) (*domain.VerificationCode, error) {
	ctx := context.Background()
	query := "SELECT id, email, code, expires_at, attempts FROM verification_codes WHERE email = $1 AND expires_at > NOW() ORDER BY created_at DESC LIMIT 1"

	var vc domain.VerificationCode
	err := r.pool.QueryRow(ctx, query, email).Scan(&vc.ID, &vc.Email, &vc.Code, &vc.ExpiresAt, &vc.Attempts)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &vc, nil
}

func (r *AuthRepository) IncrementAttempts(id int) error {
	ctx := context.Background()
	query := "UPDATE verification_codes SET attempts = attempts + 1 WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) DeleteCode(email string) error {
	ctx := context.Background()
	query := "DELETE FROM verification_codes WHERE email = $1"

	_, err := r.pool.Exec(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) CreateUser(user *domain.User) error {
	ctx := context.Background()
	query := "INSERT INTO users (username, email, password_hash, group_name, role, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	if user.Role == "" {
		user.Role = "student"
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	err := r.pool.QueryRow(
		ctx, query,
		user.Username, user.Email,
		user.PasswordHash, user.GroupName,
		user.Role, user.CreatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*domain.User, error) {
	ctx := context.Background()
	query := "SELECT id, username, email, password_hash, group_name, role, created_at FROM users WHERE email = $1"

	var u domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Username,
		&u.Email, &u.PasswordHash,
		&u.GroupName, &u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *AuthRepository) IsCodeRequestedRecently(email string) (bool, error) {
	ctx := context.Background()
	query := `SELECT EXISTS(
		SELECT 1 FROM verification_codes 
		WHERE email = $1 AND created_at > NOW() - INTERVAL '60 seconds'
	)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}