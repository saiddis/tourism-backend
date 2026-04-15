package domain

import "time"

type Tour struct {
	ID            int          `json:"id"`
	DestinationID int          `json:"destination_id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Price         float64      `json:"price"`
	StartDate     time.Time    `json:"start_date"`
	EndDate       time.Time    `json:"end_date"`
	Capacity      int          `json:"capacity"`
	Destination   *Destination `json:"destination,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}
