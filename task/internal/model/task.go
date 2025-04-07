package model

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	InProgress Status = iota
	Pending
	Completed
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      Status
	AssignedTo  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
