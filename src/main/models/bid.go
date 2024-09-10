package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// BidStatusType represents the status of a bid
type BidStatusType string

// BidAuthorType represents the type of the bid author
type BidAuthorType string

// Bid represents a bid
type Bid struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      BidStatusType `json:"status"`
	TenderID    uuid.UUID     `json:"tender_id"`
	AuthorType  BidAuthorType `json:"author_type"`
	AuthorID    string        `json:"author_id"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"created_at"`
}

// NewBid returns a new Bid instance
func NewBid(name, description string, status BidStatusType, tenderID uuid.UUID, authorType BidAuthorType, authorID string) *Bid {
	return &Bid{
		Name:        name,
		Description: description,
		Status:      status,
		TenderID:    tenderID,
		AuthorType:  authorType,
		AuthorID:    authorID,
		Version:     1,
		CreatedAt:   time.Now(),
	}
}

// GetID returns the ID of the bid
func (b *Bid) GetID() int {
	return b.ID
}

// GetName returns the name of the bid
func (b *Bid) GetName() string {
	return b.Name
}

// GetDescription returns the description of the bid
func (b *Bid) GetDescription() string {
	return b.Description
}

// GetStatus returns the status of the bid
func (b *Bid) GetStatus() BidStatusType {
	return b.Status
}

// GetTenderID returns the tender ID of the bid
func (b *Bid) GetTenderID() uuid.UUID {
	return b.TenderID
}

// GetAuthorType returns the author type of the bid
func (b *Bid) GetAuthorType() BidAuthorType {
	return b.AuthorType
}

// GetAuthorID returns the author ID of the bid
func (b *Bid) GetAuthorID() string {
	return b.AuthorID
}

// GetVersion returns the version of the bid
func (b *Bid) GetVersion() int {
	return b.Version
}

// GetCreatedAt returns the creation time of the bid
func (b *Bid) GetCreatedAt() time.Time {
	return b.CreatedAt
}

// Scan scans the bid from a database row
func (b *Bid) Scan(rows *sql.Rows) error {
	return rows.Scan(&b.ID, &b.Name, &b.Description, &b.Status, &b.TenderID, &b.AuthorType, &b.AuthorID, &b.Version, &b.CreatedAt)
}
