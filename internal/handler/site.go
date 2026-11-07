package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/EgorVarshutin05/website-checking/internal/domain"
	"github.com/EgorVarshutin05/website-checking/internal/service"
)

type SiteHandler struct {
	monitor *service.MonitorService
}

func NewSiteHandler(monitor *service.MonitorService) *SiteHandler {
	return &SiteHandler{
		monitor: monitor,
	}
}

func (h *SiteHandler) GetSites(w http.ResponseWriter, r *http.Request) {

	sites, err := h.monitor.GetSites(r.Context())
	if err != nil {
		log.Printf("Ошибка получения сайтов: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sites); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
	}

}

func (h *SiteHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSiteRequest

	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Ожидался один JSON-объект", http.StatusBadRequest)
		return
	}

	site, err := h.monitor.CreateSite(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSite):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, service.ErrSiteExists):
			http.Error(w, "Сайт уже существует", http.StatusConflict)

		default:
			log.Printf("Ошибка создания сайта: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(site); err != nil {
		log.Printf("Ошибка записи JSON: %v", err)
	}
}
