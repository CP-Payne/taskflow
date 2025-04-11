package model

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID
	Username string
	Email    string
}
