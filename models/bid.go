package models

import (
	"time"

	"github.com/google/uuid"
)

// BidAuthorType represents the type of the bid author
type BidAuthorType string

const (
	BidAuthorOrganization BidAuthorType = "Organization"
	BidAuthorUser         BidAuthorType = "User"
)

// BidStatusType represents the status of a bid
type BidStatusType string

const (
	BidStatusCreated   BidStatusType = "Created"
	BidStatusPublished BidStatusType = "Published"
	BidStatusCanceled  BidStatusType = "Canceled"
)

type BidCreate struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	TenderID    uuid.UUID     `json:"tenderId"`
	AuthorType  BidAuthorType `json:"authorType"`
	AuthorID    uuid.UUID     `json:"authorID"`
}

// Bid represents a bid
type Bid struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      BidStatusType `json:"status"`
	TenderID    uuid.UUID     `json:"tenderId"`
	AuthorType  BidAuthorType `json:"author_type"`
	AuthorID    uuid.UUID     `json:"authorId"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"createdAt"`
}

type BidUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Review struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
