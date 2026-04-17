package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type TourHighlightHandler struct {
	service *service.TourHighlightService
}

func NewTourHighlightHandler(service *service.TourHighlightService) *TourHighlightHandler {
	return &TourHighlightHandler{service: service}
}

func (h *TourHighlightHandler) Create(w http.ResponseWriter, r *http.Request) {
	tourID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tour id")
		return
	}

	var highlight domain.TourHighlight
	if err := json.NewDecoder(r.Body).Decode(&highlight); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	highlight.TourID = tourID

	if err := h.service.Create(r.Context(), &highlight); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, highlight)
}

func (h *TourHighlightHandler) GetByTourID(w http.ResponseWriter, r *http.Request) {
	tourID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tour id")
		return
	}

	highlights, err := h.service.GetByTourID(r.Context(), tourID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, highlights)
}

func (h *TourHighlightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "highlightId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid highlight id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "highlight deleted"})
}
