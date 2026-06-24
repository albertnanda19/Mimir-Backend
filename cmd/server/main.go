package main

import (
	"log"
	"net/http"
	"os"

	"github.com/albertnanda19/mimir-backend/internal/config"
	"github.com/albertnanda19/mimir-backend/internal/server"
)

func main() {
	cfg := config.Load()
	r := server.New(cfg)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	log.Printf("mimir-backend starting on :%s", addr)
	log.Fatal(http.ListenAndServe(":"+addr, r))
}
