package models

import (
	"time"
)

// Organization represents an organization
type Organization struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrganizationResponsible represents an organization responsible
type OrganizationResponsible struct {
	ID             int `json:"id"`
	OrganizationID int `json:"organization_id"`
	UserID         int `json:"user_id"`
}
