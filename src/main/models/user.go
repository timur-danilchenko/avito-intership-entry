package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents a user
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser returns a new User instance
func NewUser(username, firstName, lastName string) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetID returns the ID of the user
func (u *User) GetID() uuid.UUID {
	return u.ID
}

// GetUsername returns the username of the user
func (u *User) GetUsername() string {
	return u.Username
}

// GetFirstName returns the first name of the user
func (u *User) GetFirstName() string {
	return u.FirstName
}

// GetLastName returns the last name of the user
func (u *User) GetLastName() string {
	return u.LastName
}

// GetCreatedAt returns the creation time of the user
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// GetUpdatedAt returns the update time of the user
func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

// Scan scans the user from a database row
func (u *User) Scan(rows *sql.Rows) error {
	return rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt)
}
