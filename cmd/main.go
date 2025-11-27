package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/et0/links-checker-25-11-2025/internal/config"
	"github.com/et0/links-checker-25-11-2025/internal/handler"
	"github.com/et0/links-checker-25-11-2025/internal/logging"
	"github.com/et0/links-checker-25-11-2025/internal/repository/memory"
)

func main() {
	log := logging.New()

	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	cfg, err := config.Loading()
	if err != nil {
		log.Error("failed config load", "error", err)
		return
	}

	db, err := memory.NewLink(cfg.Storage.Filepath)
	if err != nil {
		log.Error("failed init db", "error", err)
	}

	AppHandler := handler.New(log, cfg, db)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: AppHandler.Route,
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Could not listen", "error", err)
			return
		}
	}()

	// Ожидаем сигнала из канала об остановке программы
	<-exit

	log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	AppHandler.Shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Could not gracefully shutdown the server", "error", err)
	}

	log.Info("Server stopped")
}
