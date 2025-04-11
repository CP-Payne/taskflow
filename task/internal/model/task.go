package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/CP-Payne/taskflow/pkg/gen/task/v1"
)

type Status int

const (
	InProgress Status = iota
	Pending
	Completed
)

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Status      Status
	AssignedTo  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) ToProto() *api.Task {
	var assignedTo *api.UUID
	if t.AssignedTo != nil {
		assignedTo = UuidToProtoUUID(*t.AssignedTo)
	}
	return &api.Task{
		Id:          UuidToProtoUUID(t.ID),
		UserId:      UuidToProtoUUID(t.UserID),
		Title:       t.Title,
		Description: t.Description,
		Status:      api.Status(t.Status),
		AssignedTo:  assignedTo,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
}

func UuidToProtoUUID(id uuid.UUID) *api.UUID {
	return &api.UUID{
		Value: id.String(),
	}
}

func TaskListToProto(tasks []Task) []*api.Task {
	taskProtoList := make([]*api.Task, len(tasks))
	for i, t := range tasks {
		taskProtoList[i] = t.ToProto()
	}

	return taskProtoList
}
