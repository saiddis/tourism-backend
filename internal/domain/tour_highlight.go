package domain

import "time"

type TourHighlight struct {
	ID        int       `json:"id"`
	TourID    int       `json:"tour_id"`
	ImageURL  string    `json:"image_url"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
