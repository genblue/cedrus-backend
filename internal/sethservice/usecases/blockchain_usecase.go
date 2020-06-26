package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/repository"
	pkgmodel "github.com/genblue-private/cedrus-backend/pkg/domain/model"
	pkgrepository "github.com/genblue-private/cedrus-backend/pkg/domain/repository"
	"log"
	"time"
)

type BlockchainUsecase struct {
	blockchainRepository repository.BlockchainRepository
	claimRepository      pkgrepository.ClaimRepository
	emailRepository      pkgrepository.EmailRepository
	administrator        pkgmodel.EmailRecipient
}

type BlockchainUsecaseInterface interface {
	TransferCedarCoinsToAddress(to string, claimCode string) (transaction *model.Transaction, err error)
	Health() error
	FindAccountBalance() (*model.AccountBalance, error)
}

func NewBlockchainUsecase(
	blockchainRepository repository.BlockchainRepository,
	claimRepository pkgrepository.ClaimRepository,
	emailRepository pkgrepository.EmailRepository,
	administrator pkgmodel.EmailRecipient) *BlockchainUsecase {
	return &BlockchainUsecase{
		blockchainRepository: blockchainRepository,
		claimRepository:      claimRepository,
		emailRepository:      emailRepository,
		administrator:        administrator,
	}
}

func (buc *BlockchainUsecase) TransferCedarCoinsToAddress(to string, claimCode string) (transaction *model.Transaction, errTransfer error) {
	claim, errTransfer := buc.claimRepository.FindById(claimCode)
	logInfo := fmt.Sprintf("(claimCode: \"%s\")", claimCode)
	if errTransfer != nil {
		log.Println("ERROR", "no claim found", logInfo)
		log.Println("ERROR", errTransfer, logInfo)
		return nil, errors.New("no claim found")
	}

	if claim.Status != pkgmodel.ClaimStatusOpen {
		return nil, errors.New("not authorized to claim tokens")
	}

	/* TODO: How do we manage approving ?
	   log.Println("Updating the claim as approved", logInfo)
	      claim.Status = pkgmodel.ClaimStatusApproved
	      errTransfer = buc.claimRepository.UpdateById(claim)
	      if errTransfer != nil {
	          log.Println("ERROR", errTransfer, logInfo)
	          return nil , errTransfer
	      }
	*/

	log.Println("Sending the transfer request to the blockchain", logInfo)
	transaction, errTransfer = buc.blockchainRepository.TransferCedarCoins(to, claim.TreeCount)
	if errTransfer != nil {
		log.Println("ERROR", "could not transfer claim", logInfo, ":", errTransfer)
		err := notifyAdministrator(claim, errTransfer, buc)
		if err != nil {
			log.Println("ERROR", "could not notify administrator", buc.administrator, err)
		}

		return nil, errTransfer
	}

	log.Println("Successfully sent. Updating the claim as closed", logInfo)
	claim.SettlementDate = time.Now().Unix()
	claim.TransferAddress = to
	claim.Status = pkgmodel.ClaimStatusClosed
	errTransfer = buc.claimRepository.UpdateById(claim)
	if errTransfer != nil {
		log.Println("ERROR", errTransfer, logInfo)
		return nil, errTransfer
	}

	return transaction, nil
}

func notifyAdministrator(claim *pkgmodel.Claim, err error, buc *BlockchainUsecase) error {
	claimAsJSON, _ := json.Marshal(claim)
	errorEmail := fmt.Sprintf("Could not execute token transfer for claim: %s. <br/> "+
		"Error: %s", string(claimAsJSON), err)
	err = buc.emailRepository.SendEmail(
		"cedrus@genblue.io",
		"Error with a Cedar token transfer",
		buc.administrator,
		errorEmail,
		errorEmail)
	return err
}

func (buc *BlockchainUsecase) Health() error {
	err := buc.blockchainRepository.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (buc *BlockchainUsecase) FindAccountBalance() (*model.AccountBalance, error) {
	accountBalance, err := buc.blockchainRepository.AccountBalance()
	if err != nil {
		return nil, err
	}

	return accountBalance, nil
}
