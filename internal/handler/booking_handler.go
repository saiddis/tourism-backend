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

type BookingHandler struct {
	service        *service.BookingService
	tourService    *service.TourService
	userService    *service.UserService
	paymentService *service.PaymentService
}

func NewBookingHandler(
	service *service.BookingService,
	tourService *service.TourService,
	userService *service.UserService,
	paymentService *service.PaymentService,
) *BookingHandler {
	return &BookingHandler{
		service:        service,
		tourService:    tourService,
		userService:    userService,
		paymentService: paymentService,
	}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TourID int `json:"tour_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 1. Get tour
	tour, err := h.tourService.GetByID(req.TourID)
	if err != nil {
		respondError(w, http.StatusNotFound, "tour not found")
		return
	}

	// 2. Check capacity
	if tour.Capacity <= 0 {
		respondError(w, http.StatusBadRequest, "tour is fully booked")
		return
	}

	// 3. Check for existing active booking
	existingBookings, err := h.service.GetByUserID(claims.UserID)
	if err == nil {
		for _, b := range existingBookings {
			if b.TourID == req.TourID && b.Status != domain.BookingStatusCancelled {
				respondError(w, http.StatusBadRequest, "you already have a booking for this tour")
				return
			}
		}
	}

	// 4. Deduct balance
	err = h.userService.DeductBalance(claims.UserID, tour.Price)
	if err != nil {
		respondError(w, http.StatusBadRequest, "insufficient balance")
		return
	}

	// 5. Create booking (confirmed)
	booking := &domain.Booking{
		UserID:                 claims.UserID,
		TourID:                 req.TourID,
		TourName:               tour.Name,
		TourDescription:        tour.Description,
		TourPrice:              tour.Price,
		TourStartDate:          tour.StartDate,
		TourEndDate:            tour.EndDate,
		TourCapacity:           tour.Capacity,
		DestinationID:          tour.DestinationID,
		DestinationName:        tour.DestinationName,
		DestinationDescription: tour.DestinationDescription,
		DestinationImageURL:    tour.DestinationImageURL,
		Status:                 domain.BookingStatusConfirmed,
	}

	booking, err = h.service.Create(booking)
	if err != nil {
		// Rollback: refund the user
		h.userService.RefundBalance(claims.UserID, tour.Price)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 6. Create payment (paid)
	payment := &domain.Payment{
		BookingID: booking.ID,
		Amount:    tour.Price,
		Currency:  "TJS",
	}
	if _, err := h.paymentService.Create(payment, domain.PaymentStatusPaid); err != nil {
		// Non-critical: payment record failed but booking succeeded
		// In production, you'd want to handle this more carefully
	}

	// 7. Decrement tour capacity
	if err := h.tourService.DecrementCapacity(tour.ID); err != nil {
		// Non-critical: booking succeeded
	}

	respondJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, bookings)
}

func (h *BookingHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if claims.UserID != id && claims.Role != "manager" && claims.Role != "admin" {
		respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	bookings, err := h.service.GetByUserID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, bookings)
}

func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Status domain.BookingStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims := middleware.GetClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get booking to check ownership and get tour info
	booking, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "booking not found")
		return
	}

	// Non-managers can only cancel their own bookings
	if claims.Role != "manager" && claims.Role != "admin" {
		if req.Status != domain.BookingStatusCancelled {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
		if booking.UserID != claims.UserID {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	// If cancelling a confirmed booking, refund the user
	if req.Status == domain.BookingStatusCancelled && booking.Status == domain.BookingStatusConfirmed {
		_, err := h.userService.RefundBalance(booking.UserID, booking.TourPrice)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to process refund")
			return
		}
		// Increment tour capacity
		h.tourService.IncrementCapacity(booking.TourID)
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}

func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "booking deleted"})
}
