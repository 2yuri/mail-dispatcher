package domain

type Mailer interface {
	SendMail(mail Message)
}
