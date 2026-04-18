package domain

import "time"

type Tour struct {
	ID                     int              `json:"id"`
	UserID                 int              `json:"user_id,omitempty"`
	DestinationID          int              `json:"destination_id"`
	DestinationName        string           `json:"destination_name,omitempty"`
	DestinationDescription string           `json:"destination_description,omitempty"`
	DestinationImageURL    string           `json:"destination_image_url,omitempty"`
	Name                   string           `json:"name"`
	Description            string           `json:"description"`
	Price                  float64          `json:"price"`
	StartDate              time.Time        `json:"start_date"`
	EndDate                time.Time        `json:"end_date"`
	Capacity               int              `json:"capacity"`
	ProviderID             *int             `json:"provider_id,omitempty"`
	Highlights             []*TourHighlight `json:"highlights,omitempty"`
	CreatedAt              time.Time        `json:"created_at"`
}
