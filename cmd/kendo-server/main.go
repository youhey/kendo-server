package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/youhey/kendo-server/internal/api"
	"github.com/youhey/kendo-server/internal/config"
	"github.com/youhey/kendo-server/internal/db"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := runHealthcheck(os.Args[2:]); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		return
	}

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

func runHealthcheck(args []string) error {
	target := "http://localhost:8080/healthz"
	if len(args) > 0 {
		target = args[0]
	}

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return fmt.Errorf("healthcheck request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthcheck returned status %d", resp.StatusCode)
	}

	return nil
}
