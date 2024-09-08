package models

import "time"

type Organization struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrganizationResponsible struct {
	ID             int `json:"id"`
	OrganizationID int `json:"organization_id"`
	UserID         int `json:"user_id"`
}
