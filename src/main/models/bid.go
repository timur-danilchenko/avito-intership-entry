package models

import (
	"time"
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
	TenderID    int           `json:"tender_id"`
	AuthorType  BidAuthorType `json:"author_type"`
	AuthorID    string        `json:"author_id"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Review struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
