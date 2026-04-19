package domain

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID                     int              `json:"id"`
	UserID                 int              `json:"user_id"`
	TourID                 int              `json:"tour_id"`
	TourName               string           `json:"tour_name,omitempty"`
	TourDescription        string           `json:"tour_description,omitempty"`
	TourPrice              float64          `json:"tour_price,omitempty"`
	TourStartDate          time.Time        `json:"tour_start_date,omitempty"`
	TourEndDate            time.Time        `json:"tour_end_date,omitempty"`
	TourCapacity           int              `json:"tour_capacity,omitempty"`
	RemainingSpots         int              `json:"remaining_spots,omitempty"`
	DestinationID          int              `json:"destination_id,omitempty"`
	DestinationName        string           `json:"destination_name,omitempty"`
	DestinationDescription string           `json:"destination_description,omitempty"`
	DestinationImageURL    string           `json:"destination_image_url,omitempty"`
	TourHighlights         []*TourHighlight `json:"tour_highlights,omitempty"`
	Status                 BookingStatus    `json:"status"`
	CreatedAt              time.Time        `json:"created_at"`
}
