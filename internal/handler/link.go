package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/et0/links-checker-25-11-2025/internal/model"
	"github.com/et0/links-checker-25-11-2025/internal/service"
)

type Link struct {
	linkService service.LinkService
}

func NewLinkHandler(linkService service.LinkService) *Link {
	return &Link{
		linkService,
	}
}

func (l *Link) Check(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Links) == 0 {
		http.Error(w, "No links provided", http.StatusBadRequest)
		return
	} else if len(req.Links) > 10 {
		http.Error(w, "10 links limit", http.StatusBadRequest)
		return
	}

	response, err := l.linkService.Check(r.Context(), req.Links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (l *Link) Report(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(req.LinksIDs)

	if len(req.LinksIDs) == 0 {
		http.Error(w, "No links provided", http.StatusBadRequest)
		return
	} else if len(req.LinksIDs) > 10 {
		http.Error(w, "10 links limit", http.StatusBadRequest)
		return
	}

	pdf, err := l.linkService.Report(r.Context(), req.LinksIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
	w.Write(pdf)
}
