package models

import (
	"time"

	"github.com/google/uuid"
)

type TenderServiceType string

const (
	TenderServiceConstruction TenderServiceType = "Construction"
	TenderServiceDelivery     TenderServiceType = "Delivery"
	TenderServiceManufacture  TenderServiceType = "Manufacture"
)

type TenderStatusType string

const (
	TenderStatusCreated   TenderStatusType = "Created"
	TenderStatusPublished TenderStatusType = "Published"
	TenderStatusClosed    TenderStatusType = "Closed"
)

type TenderCreate struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	ServiceType     TenderServiceType `json:"serviceType"`
	OrganizationID  uuid.UUID         `json:"organizationId"`
	CreatorUsername string            `json:"creatorUsername"`
}

// Tender represents a tender
type Tender struct {
	ID             uuid.UUID         `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	ServiceType    TenderServiceType `json:"serviceType"`
	Status         TenderStatusType  `json:"status"`
	OrganizationID uuid.UUID         `json:"organizationId"`
	Version        int               `json:"version"`
	CreatedAt      time.Time         `json:"createdAt"`
}

// TenderUpdate represents a update tender message
type TenderUpdate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceType TenderServiceType `json:"serviceType"`
}
