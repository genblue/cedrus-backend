package usecases

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/genblue-private/cedrus-backend/pkg/domain/repository"
	"github.com/sirupsen/logrus"
	"log"
)

type ClaimUsecase struct {
	repo repository.ClaimRepository
}

type ClaimUsecaseInterface interface {
	SaveClaim(claim *model.Claim) error
	FindClaims() ([]*model.Claim, error)
	FindClaim(id string) (*model.Claim, error)
}

func NewClaimUsecase(repo repository.ClaimRepository) *ClaimUsecase {
	return &ClaimUsecase{
		repo: repo,
	}
}

func (cuc *ClaimUsecase) SaveClaim(claim *model.Claim) error {
	err := checkmail.ValidateFormat(claim.Email)
	if err != nil {
		log.Println("WARNING", "bad email format:", claim.Email)
		return errors.New(err.Error() + ": " + claim.Email)
	}

	/* Validate host disabled for now as it's instable
	   err = checkmail.ValidateHost(claim.Email)
	   if err != nil {
	       log.Println("WARNING", "bad email host:", claim.Email)
	       return errors.New(err.Error() + ": " + claim.Email)
	   }
	*/

	err = cuc.repo.Save(claim)
	if err != nil {
		log.Printf("Could not save claim %p", claim)
		return err
	}

	logrus.Debug("Saved claim:", claim)
	return nil
}

func (cuc *ClaimUsecase) FindClaims() ([]*model.Claim, error) {
	claims, err := cuc.repo.FindAll()

	if err != nil {
		log.Printf("Could not find claims: %p", err)
		return nil, err
	}

	logrus.Debug("Find claims:", claims)
	return claims, nil
}

func (cuc *ClaimUsecase) FindClaim(id string) (*model.Claim, error) {
	claim, err := cuc.repo.FindById(id)

	if err != nil {
		log.Printf("Could not find claim: %p", err)
		return nil, err
	}

	logrus.Debug("Find claim:", claim)
	return claim, nil
}
