package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization
type Organization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewOrganization returns a new Organization instance
func NewOrganization(name, description, orgType string) *Organization {
	return &Organization{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Type:        orgType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// GetID returns the ID of the organization
func (o *Organization) GetID() uuid.UUID {
	return o.ID
}

// GetName returns the name of the organization
func (o *Organization) GetName() string {
	return o.Name
}

// GetDescription returns the description of the organization
func (o *Organization) GetDescription() string {
	return o.Description
}

// GetType returns the type of the organization
func (o *Organization) GetType() string {
	return o.Type
}

// GetCreatedAt returns the creation time of the organization
func (o *Organization) GetCreatedAt() time.Time {
	return o.CreatedAt
}

// GetUpdatedAt returns the update time of the organization
func (o *Organization) GetUpdatedAt() time.Time {
	return o.UpdatedAt
}

// Scan scans the organization from a database row
func (o *Organization) Scan(rows *sql.Rows) error {
	return rows.Scan(&o.ID, &o.Name, &o.Description, &o.Type, &o.CreatedAt, &o.UpdatedAt)
}

// OrganizationResponsible represents an organization responsible
type OrganizationResponsible struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
}

// NewOrganizationResponsible returns a new OrganizationResponsible instance
func NewOrganizationResponsible(organizationID, userID uuid.UUID) *OrganizationResponsible {
	return &OrganizationResponsible{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		UserID:         userID,
	}
}

// GetID returns the ID of the organization responsible
func (or *OrganizationResponsible) GetID() uuid.UUID {
	return or.ID
}

// GetOrganizationID returns the organization ID of the organization responsible
func (or *OrganizationResponsible) GetOrganizationID() uuid.UUID {
	return or.OrganizationID
}

// GetUserID returns the user ID of the organization responsible
func (or *OrganizationResponsible) GetUserID() uuid.UUID {
	return or.UserID
}

// Scan scans the organization responsible from a database row
func (or *OrganizationResponsible) Scan(rows *sql.Rows) error {
	return rows.Scan(&or.ID, &or.OrganizationID, &or.UserID)
}
