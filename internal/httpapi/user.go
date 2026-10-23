package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MishaEn/techstore-user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
	}

	if err := a.UserRepo.Create(r.Context(), user); err != nil {
		// Выводим реальную ошибку в консоль сервера
		slog.Error("failed to create user in DB", "error", err)

		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
