package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
	middleware "tourism-backend/internal/middlewares"
)

type TourHandler struct {
	service         *service.TourService
	providerService *service.ProviderService
}

func NewTourHandler(service *service.TourService, providerService *service.ProviderService) *TourHandler {
	return &TourHandler{service: service, providerService: providerService}
}

func (h *TourHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tour domain.Tour
	if err := json.NewDecoder(r.Body).Decode(&tour); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := h.service.Create(r.Context(), &tour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *TourHandler) CreateForProvider(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	provider, err := h.providerService.GetByUserID(r.Context(), claims.UserID)
	if err != nil || provider == nil || !provider.Active {
		respondError(w, http.StatusForbidden, "no active provider profile")
		return
	}

	var tour domain.Tour
	if err := json.NewDecoder(r.Body).Decode(&tour); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tour.ProviderID = &provider.ID

	result, err := h.service.Create(r.Context(), &tour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *TourHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tours, err := h.service.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.service.CalculateRemainingSpots(r.Context(), tours); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tours)
}

func (h *TourHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tour, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	remaining, err := h.service.GetRemainingSpots(r.Context(), tour.ID)
	if err == nil {
		tour.RemainingSpots = remaining
	}
	respondJSON(w, http.StatusOK, tour)
}

func (h *TourHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var tour domain.Tour
	if err := json.NewDecoder(r.Body).Decode(&tour); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	tour.ID = id
	result, err := h.service.Update(r.Context(), &tour)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *TourHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "tour deleted"})
}
func (h *TourHandler) GetByDestinationID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tours, err := h.service.GetByDestinationID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.service.CalculateRemainingSpots(r.Context(), tours); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tours)
}

func (h *TourHandler) Renew(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	provider, err := h.providerService.GetByUserID(r.Context(), claims.UserID)
	if err != nil || provider == nil || !provider.Active {
		respondError(w, http.StatusForbidden, "no active provider profile")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	tour, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "tour not found")
		return
	}

	if tour.ProviderID == nil || *tour.ProviderID != provider.ID {
		respondError(w, http.StatusForbidden, "not authorized to renew this tour")
		return
	}

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StartDate == "" || req.EndDate == "" {
		respondError(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	if err := h.service.RenewTour(r.Context(), id, req.StartDate, req.EndDate); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "tour renewed successfully"})
}
