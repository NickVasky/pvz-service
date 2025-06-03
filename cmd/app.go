package main

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/config"
	"AvitoTechPVZ/repo"
	"AvitoTechPVZ/security"
	"AvitoTechPVZ/service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"log"
	"net/http"
)

const (
	shutdownPeriod = 10 * time.Second
)

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error loading App config: %v", err)
	}

	db, err := repo.NewDbConn(cfg.Db)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}

	server := BuildServer(db, cfg)

	go func() {
		log.Printf("Starting server on %v...", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	<-rootCtx.Done()
	Shutdown(server, db)
}

func BuildServer(db *sql.DB, cfg config.ServiceConfig) *http.Server {
	repo := repo.NewRepo(db)
	sc := security.NewSecurityController(cfg.App.Secret)
	s := service.NewDefaultAPIServicerImpl(repo, sc)

	controller := dto.NewDefaultAPIController(s)
	router := s.NewRouter(controller)

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.App.Port),
		Handler: router,
	}
}

func Shutdown(server *http.Server, db *sql.DB) {
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}

	repo.CloseDbConn(db)

	log.Println("Graceful shutdown complete.")
}
