package models

import "time"

// User represents an employee in the database.
type User struct {
	ID        int       `json:"id"`         // SERIAL PRIMARY KEY
	Username  string    `json:"username"`   // VARCHAR(50) UNIQUE NOT NULL
	FirstName string    `json:"first_name"` // VARCHAR(50)
	LastName  string    `json:"last_name"`  // VARCHAR(50)
	CreatedAt time.Time `json:"created_at"` // TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	UpdatedAt time.Time `json:"updated_at"` // TIMESTAMP DEFAULT CURRENT_TIMESTAMP
}
