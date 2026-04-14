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

type PaymentHandler struct {
	service        *service.PaymentService
	bookingService *service.BookingService
}

func NewPaymentHandler(service *service.PaymentService, bookingService *service.BookingService) *PaymentHandler {
	return &PaymentHandler{
		service:        service,
		bookingService: bookingService,
	}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payment domain.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	booking, err := h.bookingService.GetByID(payment.BookingID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if claims.Role != "manager" && claims.Role != "admin" && booking.UserID != claims.UserID {
		respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	result, err := h.service.Create(&payment)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	payment, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Status domain.PaymentStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "payment status updated"})
}
