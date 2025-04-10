package notification

import (
	"context"

	"go.uber.org/zap"
)

type Sender interface {
	Send(ctx context.Context, recipient string, message string) error
}

type LogSender struct {
	logger *zap.SugaredLogger
}

func NewLogSender(logger *zap.SugaredLogger) *LogSender {
	return &LogSender{
		logger: logger,
	}
}

func (s *LogSender) Send(ctx context.Context, recipient string, message string) error {
	s.logger.Infow("Sending notification", "recipient", recipient, "message", message)
	// Implement Email call here
	return nil
}
