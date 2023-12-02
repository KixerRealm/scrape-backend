// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type BlogPost struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Description   string
	ImageFilename string
	UserID        uuid.UUID
}

type BugReport struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Description   string
	ImageFilename string
	UserID        uuid.UUID
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string
	Username  string
	Token     string
}
