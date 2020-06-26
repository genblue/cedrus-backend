package email

import (
	"errors"
	"fmt"
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/genblue-private/cedrus-backend/pkg/domain/repository"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type SendgridEmailRepository struct {
	SendgridApiKey string
}

func NewSendgridEmailRepository(sendgridApiKey string) repository.EmailRepository {
	return &SendgridEmailRepository{
		SendgridApiKey: sendgridApiKey,
	}
}

func (smr *SendgridEmailRepository) SendEmail(
	fromAddress string,
	subject string,
	recipient model.EmailRecipient,
	plaintextBody string,
	htmlBody string) error {

	from := mail.NewEmail("CedarCoin", fromAddress)
	to := mail.NewEmail(recipient.Name, recipient.Address)
	message := mail.NewSingleEmail(from, subject, to, plaintextBody, htmlBody)
	client := sendgrid.NewSendClient(smr.SendgridApiKey)

	response, err := client.Send(message)
	if err != nil {
		log.Println("Error while sending an email", err)
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		log.Println("Error while sending an email", response)
		errorMessage := fmt.Sprintf("Bad API status code from Sendgrid: %v", response.StatusCode)
		return errors.New(errorMessage)
	}

	logrus.Debug("Sent an email:", message)
	return nil
}
