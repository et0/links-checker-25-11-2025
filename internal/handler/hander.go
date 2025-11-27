package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/et0/links-checker-25-11-2025/internal/config"
	"github.com/et0/links-checker-25-11-2025/internal/repository"
	"github.com/et0/links-checker-25-11-2025/internal/service"
)

type AppHandler struct {
	linkService service.LinkService
	Route       *http.ServeMux
}

func New(log *slog.Logger, cfg *config.Config, db repository.LinkRepository) *AppHandler {
	mux := http.NewServeMux()

	// Service
	linkService := service.NewLinkService(db, cfg)

	// Handler
	linkHandler := NewLinkHandler(linkService)

	go linkHandler.linkService.CheckContinue(context.Background())

	mux.HandleFunc("/check", linkHandler.Check)
	mux.HandleFunc("/report", linkHandler.Report)
	// mux.HandleFunc("/generate-report", httpHandler.GenerateReport)
	// mux.HandleFunc("/batch-status", httpHandler.GetBatchStatus)

	return &AppHandler{
		linkService,
		mux,
	}
}

func (h *AppHandler) Shutdown() {
	h.linkService.Shutdown()
}
