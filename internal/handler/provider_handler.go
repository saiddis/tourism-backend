package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tourism-backend/internal/domain"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type ProviderHandler struct {
	service *service.ProviderService
}

func NewProviderHandler(service *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var provider domain.Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), &provider)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *ProviderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	providers, err := h.service.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, providers)
}

func (h *ProviderHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	providers, err := h.service.GetActiveProviders(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, providers)
}

func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	provider, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, provider)
}

func (h *ProviderHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	provider, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if provider == nil {
		respondError(w, http.StatusNotFound, "provider not found")
		return
	}
	respondJSON(w, http.StatusOK, provider)
}

func (h *ProviderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var provider domain.Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	provider.ID = id

	result, err := h.service.Update(r.Context(), &provider)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *ProviderHandler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ToggleActive(r.Context(), id, req.Active); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "provider active status updated"})
}
