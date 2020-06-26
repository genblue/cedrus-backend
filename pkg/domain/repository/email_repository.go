package repository

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
)

type EmailRepository interface {
	SendEmail(
		from string,
		subject string,
		recipient model.EmailRecipient,
		plaintextBody string,
		htmlBody string) error
}
