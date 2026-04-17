package domain

import "time"

type ProviderType string

const (
	ProviderTypePrivate      ProviderType = "private"
	ProviderTypeOrganization ProviderType = "organization"
)

type Provider struct {
	ID              int          `json:"id"`
	UserID          int          `json:"user_id"`
	UserName        string       `json:"user_name,omitempty"`
	UserEmail       string       `json:"user_email,omitempty"`
	UserAvatarURL   *string      `json:"user_avatar_url,omitempty"`
	Phone           string       `json:"phone"`
	ProviderType    ProviderType `json:"provider_type"`
	InstagramURL    *string      `json:"instagram_url,omitempty"`
	TelegramURL     *string      `json:"telegram_url,omitempty"`
	FacebookURL     *string      `json:"facebook_url,omitempty"`
	YearsExperience int          `json:"years_experience"`
	Bio             *string      `json:"bio,omitempty"`
	Active          bool         `json:"active"`
	CreatedAt       time.Time    `json:"created_at"`
}
