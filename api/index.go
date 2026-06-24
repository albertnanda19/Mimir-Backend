package handler

import (
	"net/http"

	"github.com/albertnanda19/mimir-backend/internal/config"
	"github.com/albertnanda19/mimir-backend/internal/server"
)

var router http.Handler

func init() {
	cfg := config.Load()
	router = server.New(cfg)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
