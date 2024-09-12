package models

import (
	"time"

	"github.com/google/uuid"
)

type TenderCreate struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"` // enum
	OrganizationID  uuid.UUID `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
}

// Tender represents a tender
type Tender struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"serviceType"` // enum
	Status         string    `json:"status"`      // enum
	OrganizationID uuid.UUID `json:"organizationId"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
}

// TenderUpdate represents a update tender message
type TenderUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}
