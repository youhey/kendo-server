package main

import (
	"log"
	"net/http"

	"github.com/youhey/kendo-server/internal/api"
	"github.com/youhey/kendo-server/internal/config"
	"github.com/youhey/kendo-server/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	store, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer store.Close()

	handler := api.NewHandler(store)
	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	log.Printf("kendo-server listening on %s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
