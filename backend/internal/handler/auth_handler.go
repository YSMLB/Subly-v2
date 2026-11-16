package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"subly-v2/internal/domain"
	"subly-v2/internal/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

type SendCodeRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponseUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	GroupName string `json:"group_name"`
}

type RegisterResponse struct {
	Token string               `json:"token"`
	User  RegisterResponseUser `json:"user"`
}

func (h *AuthHandler) SendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	if req.Email == "" || !strings.Contains(req.Email, "@") {
		http.Error(w, "некорректный email", http.StatusBadRequest)
		return
	}

	err = h.AuthService.RequestCode(req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrRateLimited) {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "verification code sent",
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Code == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "все обязательные поля должны быть заполнены", http.StatusBadRequest)
		return
	}

	token, user, err := h.AuthService.Register(req)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, domain.ErrInvalidCode) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := RegisterResponse{
		Token: token,
		User: RegisterResponseUser{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			GroupName: user.GroupName,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "некорректное тело запроса", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email и пароль обязательны", http.StatusBadRequest)
		return
	}

	token, user, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	res := RegisterResponse{
		Token: token,
		User: RegisterResponseUser{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			GroupName: user.GroupName,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
