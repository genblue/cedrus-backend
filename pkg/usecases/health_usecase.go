package usecases

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/genblue-private/cedrus-backend/pkg/domain/repository"
	"github.com/sirupsen/logrus"
	"log"
)

type HealthUsecase struct {
	repo repository.ClaimRepository
}

type HealthUsecaseInterface interface {
	Health() *model.Health
}

func NewHealthUsecase(repo repository.ClaimRepository) *HealthUsecase {
	return &HealthUsecase{
		repo: repo,
	}
}

func (cuc *HealthUsecase) Health() *model.Health {
	err := cuc.repo.Ping()

	if err != nil {
		log.Printf("Could not ping DB: %s", err)
		health := model.NewHealth(false)
		return &health
	}

	logrus.Debug("Successfully pinged DB")
	health := model.NewHealth(true)
	return &health
}
