package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" означает, что это поле не попадет в JSON ответа
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// Входные данные для регистрации
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// Входные данные для логина
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
