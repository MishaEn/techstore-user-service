package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	DB *pgxpool.Pool
}

func NewRouter(db *pgxpool.Pool) http.Handler {
	api := &API{DB: db}
	r := chi.NewRouter()

	// Базовые middleware (логирование запросов, восстановление после паник)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Health checks
	r.Get("/healthz", api.handleHealthz)
	r.Get("/readyz", api.handleReadyz)

	// Группируем API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Сюда мы добавим /users и /auth на следующем шаге
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
