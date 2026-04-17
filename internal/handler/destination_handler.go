package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type DestinationHandler struct {
	service *service.DestinationService
}

func NewDestinationHandler(service *service.DestinationService) *DestinationHandler {
	return &DestinationHandler{service: service}
}

func (h *DestinationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var destination domain.Destination
	if err := json.NewDecoder(r.Body).Decode(&destination); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := h.service.Create(r.Context(), &destination)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *DestinationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	destinations, err := h.service.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, destinations)
}

func (h *DestinationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	destination, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, destination)
}

func (h *DestinationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "destination deleted"})
}
