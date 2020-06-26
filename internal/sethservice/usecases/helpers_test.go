package usecases

import (
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	pkgmodel "github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/stretchr/testify/mock"
)

type mockBlockchainRepository struct {
	mock.Mock
}

func (br *mockBlockchainRepository) TransferCedarCoins(to string, amount uint) (transaction *model.Transaction, err error) {
	args := br.Called(to, amount)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (br *mockBlockchainRepository) Ping() error {
	args := br.Called()
	return args.Error(0)
}

func (br *mockBlockchainRepository) AccountBalance() (*model.AccountBalance, error) {
	args := br.Called()
	return args.Get(0).(*model.AccountBalance), nil
}

type mockClaimRepository struct {
	mock.Mock
}

func (cr mockClaimRepository) FindAll() ([]*pkgmodel.Claim, error) {
	args := cr.Called()
	return args.Get(0).([]*pkgmodel.Claim), args.Error(1)
}

func (cr mockClaimRepository) FindById(id string) (*pkgmodel.Claim, error) {
	args := cr.Called(id)
	return args.Get(0).(*pkgmodel.Claim), nil
}

func (cr mockClaimRepository) Save(claim *pkgmodel.Claim) error {
	args := cr.Called(claim)
	return args.Error(0)
}

func (cr mockClaimRepository) UpdateById(claim *pkgmodel.Claim) error {
	args := cr.Called(claim)
	return args.Error(0)
}

func (cr mockClaimRepository) FindAllByEmailUnsent() ([]*pkgmodel.Claim, error) {
	args := cr.Called()
	return args.Get(0).([]*pkgmodel.Claim), nil
}

func (cr mockClaimRepository) Ping() error {
	args := cr.Called()
	return args.Error(0)
}

type mockEmailRepository struct {
	mock.Mock
}

func (er mockEmailRepository) SendEmail(
	from string,
	subject string,
	recipient pkgmodel.EmailRecipient,
	plaintextBody string,
	htmlBody string) error {
	args := er.Called(from, subject, recipient, plaintextBody, htmlBody)
	return args.Error(0)
}
