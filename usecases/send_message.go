package usecases

import (
	"email-dispatcher/domain"
)

type emailSender struct {
	mailer domain.Mailer
}

func NewEmailSender(mailer domain.Mailer) *emailSender {
	return &emailSender{
		mailer: mailer,
	}
}

func (e *emailSender) SendEmail(email domain.Message) {
	e.mailer.SendMail(email)
}
