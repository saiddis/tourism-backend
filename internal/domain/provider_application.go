package domain

import "time"

type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusAccepted ApplicationStatus = "accepted"
	ApplicationStatusRejected ApplicationStatus = "rejected"
)

type ProviderApplication struct {
	ID              int               `json:"id"`
	UserID          int               `json:"user_id"`
	UserName        string            `json:"user_name,omitempty"`
	UserEmail       string            `json:"user_email,omitempty"`
	Phone           string            `json:"phone"`
	ProviderType    ProviderType      `json:"provider_type"`
	InstagramURL    *string           `json:"instagram_url,omitempty"`
	TelegramURL     *string           `json:"telegram_url,omitempty"`
	FacebookURL     *string           `json:"facebook_url,omitempty"`
	YearsExperience int               `json:"years_experience"`
	Bio             *string           `json:"bio,omitempty"`
	Status          ApplicationStatus `json:"status"`
	AdminNote       *string           `json:"admin_note,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
