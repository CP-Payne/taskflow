package notification

import (
	"context"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	source string
	logger *zap.SugaredLogger
	dialer *gomail.Dialer
}

func NewEmailSender(sourceEmail, appPass string, logger *zap.SugaredLogger) *EmailSender {
	dialer := gomail.NewDialer("smtp.gmail.com", 587, sourceEmail, appPass)
	return &EmailSender{
		source: sourceEmail,
		dialer: dialer,
		logger: logger,
	}
}

func (s *EmailSender) Send(ctx context.Context, recipient string, message string) error {
	messageStructure := gomail.NewMessage()

	messageStructure.SetHeader("From", s.source)
	messageStructure.SetHeader("To", recipient)
	messageStructure.SetHeader("Subject", "New Task Assigned")

	messageStructure.SetBody("text/plain", message)

	if err := s.dialer.DialAndSend(messageStructure); err != nil {
		s.logger.Panicw("Failed sending notification over gmail", "error", err, "recipient", recipient)
	} else {
		s.logger.Infow("Email sent successfully", "recipient", recipient)
	}
	return nil
}
