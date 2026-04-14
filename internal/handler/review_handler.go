package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tourism-backend/internal/domain"
	middleware "tourism-backend/internal/middlewares"
	"tourism-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	var review domain.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	review.UserID = claims.UserID
	result, err := h.service.Create(&review)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *ReviewHandler) GetByTourID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	reviews, err := h.service.GetByTourID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "review deleted"})
}
