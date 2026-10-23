package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MishaEn/techstore-user-service/internal/auth"
	"github.com/MishaEn/techstore-user-service/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	DB           *pgxpool.Pool
	UserRepo     *repository.UserRepository
	TokenManager *auth.TokenManager
}

func NewRouter(db *pgxpool.Pool, tm *auth.TokenManager) http.Handler {
	api := &API{
		DB:           db,
		UserRepo:     repository.NewUserRepository(db),
		TokenManager: tm,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", api.handleHealthz)
	r.Get("/readyz", api.handleReadyz)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users/register", api.handleRegister)
		r.Post("/auth/login", api.handleLogin)

		// Защищенные маршруты
		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware)
			r.Get("/users/me", api.handleGetMe)
		})
	})

	return r
}

func (a *API) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (a *API) handleReadyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := a.DB.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "database unreachable"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
