package infra

import (
	"bytes"
	"email-dispatcher/config"
	"email-dispatcher/domain"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type amazonSES struct {
	config *aws.Config
}

func NewAmazonSES(conf *aws.Config) *amazonSES {
	return &amazonSES{
		config: conf,
	}
}

func (a *amazonSES) SendMail(m domain.Message) {
	session := session.Must(session.NewSession(a.config))
	svc := ses.New(session)

	message, err := a.createMessage(m)
	if err != nil {
		log.Print("Error when creating message: ", err)
	}

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: message,
		},
	}

	result, err := svc.SendRawEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	log.Print(result)
}

func (a *amazonSES) createMessage(m domain.Message) ([]byte, error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("To", m.Recipient)
	msg.SetHeader("CC", m.CCRecipient()...)
	msg.SetHeader("BCC", m.BCCRecipient()...)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.HTMLBody)

	if len(m.Attachments) > 0 {
		id := uuid.New().String()
		err := a.createTemporaryFolder(id)
		if err != nil {
			return nil, err
		}
		defer a.removeTemporaryFolder(fmt.Sprintf("%v/temp/%v", config.GetConfig().RootPath(), id))
		for _, value := range m.Attachments {
			a.decodeBase64(id, value.Data, value.Name)
			msg.Attach(fmt.Sprintf("%v/temp/%v/%v", config.GetConfig().RootPath(), id, value.Name))
		}
	}

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	return emailRaw.Bytes(), nil
}

func (a *amazonSES) createTemporaryFolder(fileName string) error {
	return os.Mkdir(fmt.Sprintf("%v/temp/%v", config.GetConfig().RootPath(), fileName), 0755)
}

func (a *amazonSES) removeTemporaryFolder(fileName string) error {
	return os.RemoveAll(fileName)
}

func (a *amazonSES) decodeBase64(id string, raw string, file string) error {
	dec, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%v/temp/%v/%v", config.GetConfig().RootPath(), id, file))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}
