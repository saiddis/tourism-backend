package domain

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	TourID      int           `json:"tour_id"`
	Status      BookingStatus `json:"status"`
	Tour        *Tour         `json:"tour,omitempty"`
	Destination *Destination  `json:"destination,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}
