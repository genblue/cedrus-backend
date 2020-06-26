package usecases

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	pkgrepository "github.com/genblue-private/cedrus-backend/pkg/domain/repository"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type EmailUsecase struct {
	emailRepository pkgrepository.EmailRepository
	claimRepository pkgrepository.ClaimRepository
}

type EmailUsecaseInterface interface {
	SendEmailsToNewClaims() error
}

func NewEmailUsecase(claimRepository pkgrepository.ClaimRepository, emailRepository pkgrepository.EmailRepository) *EmailUsecase {
	return &EmailUsecase{
		emailRepository: emailRepository,
		claimRepository: claimRepository,
	}
}

func (euc *EmailUsecase) SendEmailsToNewClaims() error {
	log.Println("fetching unsent emails...")
	claims, err := euc.claimRepository.FindAllByEmailUnsent()
	if err != nil {
		return err
	}

	if len(claims) == 0 {
		log.Println("No email to send, aborting")
		return nil
	}

	log.Println("Found", len(claims), "new claims, sending email...")
	var emailsSent int
	for _, claim := range claims {
		err := sendEmail(claim, euc)
		if err != nil {
			log.Println("ERROR", "while sending email for email address", claim.Email, ":", err)
			continue
		}

		updateClaimAsEmailSent(claim)
		err = euc.claimRepository.UpdateById(claim)
		emailsSent += 1
	}

	log.Println("Sent", emailsSent, "emails")
	return nil
}

func sendEmail(claim *model.Claim, euc *EmailUsecase) error {
	err := checkmail.ValidateFormat(claim.Email)
	if err != nil {
		return errors.New(err.Error() + ": " + claim.Email)
	}

	/* Validate host disabled for now as it's instable
	   err = checkmail.ValidateHost(claim.Email)
	   if err != nil {
	       return errors.New(err.Error() + ": " + claim.Email)
	   }
	*/

	body, err := ioutil.ReadFile("email/email-plaintext.txt")
	if err != nil {
		return errors.New("Could not read email template : " + err.Error())
	}
	bodyWithClaimCode := strings.Replace(string(body), "{CLAIM_CODE}", claim.ClaimCode, -1)

	recipient := model.EmailRecipient{
		Address: claim.Email,
		Name:    claim.Name,
	}
	err = euc.emailRepository.SendEmail(
		"cedrus@genblue.io",
		"Collect your CedarCoins",
		recipient,
		bodyWithClaimCode,
		bodyWithClaimCode)
	return err
}

func updateClaimAsEmailSent(claim *model.Claim) {
	claim.Status = model.ClaimStatusOpen
	claim.EmailSentDate = time.Now().Unix()
	claim.EmailSent = true
}
