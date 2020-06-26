package usecases

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/stretchr/testify/mock"
)

type mockClaimRepository struct {
	mock.Mock
}

func (cr mockClaimRepository) FindAll() ([]*model.Claim, error) {
	args := cr.Called()
	return args.Get(0).([]*model.Claim), args.Error(1)
}

func (cr mockClaimRepository) FindById(id string) (*model.Claim, error) {
	args := cr.Called(id)
	return args.Get(0).(*model.Claim), args.Error(1)
}

func (cr mockClaimRepository) Save(claim *model.Claim) error {
	args := cr.Called(claim)
	return args.Error(0)
}

func (cr mockClaimRepository) UpdateById(claim *model.Claim) error {
	args := cr.Called(claim)
	return args.Error(0)
}

func (cr mockClaimRepository) FindAllByEmailUnsent() ([]*model.Claim, error) {
	args := cr.Called()
	return args.Get(0).([]*model.Claim), nil
}

func (cr mockClaimRepository) Ping() error {
	args := cr.Called()
	return args.Error(0)
}
