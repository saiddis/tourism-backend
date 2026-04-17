package domain

import "time"

type UserRole string

const (
	RoleClient  UserRole = "client"
	RoleManager UserRole = "manager"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // json:"-" означает: не показывать в JSON ответе
	Role         UserRole  `json:"role"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
}
