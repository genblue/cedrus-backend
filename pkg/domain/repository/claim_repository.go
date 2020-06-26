package repository

import "github.com/genblue-private/cedrus-backend/pkg/domain/model"

type ClaimRepository interface {
	FindAll() ([]*model.Claim, error)
	FindById(id string) (*model.Claim, error)
	Save(claim *model.Claim) error
	UpdateById(claim *model.Claim) error
	FindAllByEmailUnsent() ([]*model.Claim, error)
	Ping() error
}
