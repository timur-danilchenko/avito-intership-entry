package models

import (
	"time"

	"github.com/google/uuid"
)

// BidStatusType represents the status of a bid
type BidStatusType string

// BidAuthorType represents the type of the bid author
type BidAuthorType string

// Bid represents a bid
type Bid struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      BidStatusType `json:"status"`
	TenderID    uuid.UUID     `json:"tender_id"`
	AuthorType  BidAuthorType `json:"author_type"`
	AuthorID    uuid.UUID     `json:"author_id"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Review struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
