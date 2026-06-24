package server

import (
	"encoding/json"
	"net/http"

	"github.com/albertnanda19/mimir-backend/internal/config"
	"github.com/albertnanda19/mimir-backend/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.Logging)

	// ponytail: routes registered here as handlers are built
	r.Get("/", health)
	r.Get("/api/health", health)

	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
