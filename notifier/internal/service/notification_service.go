package service

import (
	"context"
	"fmt"

	"github.com/CP-Payne/taskflow/notifier/internal/gateway/user"
	"github.com/CP-Payne/taskflow/notifier/internal/notification"
	"github.com/google/uuid"
)

// TODO: Turn userGateway into interface
type NotificationService struct {
	userGateway *user.Gateway
	emailSender notification.Sender
}

func NewNotificationService(userGateway *user.Gateway, sender notification.Sender) *NotificationService {
	return &NotificationService{
		userGateway: userGateway,
		emailSender: sender,
	}
}

func (s *NotificationService) NotifyUserToCompleteTask(ctx context.Context, userID, taskID uuid.UUID) error {
	user, err := s.userGateway.GetUserDetails(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// TODO: Update Redis event to include additional task information, then updated task message
	msg := fmt.Sprintf("Hi %s, you have a task (%s) to complete!", user.Username, taskID)
	return s.emailSender.Send(ctx, user.Email, msg)
}
