package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Tender represents a tender
type Tender struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"serviceType"`
	Status         string    `json:"status"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
}

// NewTender returns a new Tender instance
func NewTender(name, description, serviceType, status string, organizationID uuid.UUID) *Tender {
	return &Tender{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		ServiceType:    serviceType,
		Status:         status,
		OrganizationID: organizationID,
		Version:        1,
		CreatedAt:      time.Now(),
	}
}

// GetID returns the ID of the tender
func (t *Tender) GetID() uuid.UUID {
	return t.ID
}

// GetName returns the name of the tender
func (t *Tender) GetName() string {
	return t.Name
}

// GetDescription returns the description of the tender
func (t *Tender) GetDescription() string {
	return t.Description
}

// GetServiceType returns the service type of the tender
func (t *Tender) GetServiceType() string {
	return t.ServiceType
}

// GetStatus returns the status of the tender
func (t *Tender) GetStatus() string {
	return t.Status
}

// GetOrganizationID returns the organization ID of the tender
func (t *Tender) GetOrganizationID() uuid.UUID {
	return t.OrganizationID
}

// GetVersion returns the version of the tender
func (t *Tender) GetVersion() int {
	return t.Version
}

// GetCreatedAt returns the creation time of the tender
func (t *Tender) GetCreatedAt() time.Time {
	return t.CreatedAt
}

// Scan scans the tender from a database row
func (t *Tender) Scan(rows *sql.Rows) error {
	return rows.Scan(&t.ID, &t.Name, &t.Description, &t.ServiceType, &t.Status, &t.OrganizationID, &t.Version, &t.CreatedAt)
}
