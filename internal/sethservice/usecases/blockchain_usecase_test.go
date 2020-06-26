package usecases

import (
	"errors"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	pkgmodel "github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"
	"testing"
)

func TestBlockchainUsecase_TransferCedarCoinsToAddress_WithValidClaim(t *testing.T) {
	// Given
	dummyClaim := pkgmodel.NewClaim("Mike", "email@dummy.fr", 45)
	dummyClaim.Status = pkgmodel.ClaimStatusOpen
	expectedTx := model.Transaction{
		TxID:        "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
		BlockNumber: "1234",
	}
	blockchainRepository := stubBlockchainRepository(&expectedTx, nil)
	claimRepository := stubClaimRepository(&dummyClaim)
	administrator := pkgmodel.EmailRecipient{
		Address: "dummy@email.com",
		Name:    "Dummy Name",
	}
	blockchainUsecase := NewBlockchainUsecase(blockchainRepository, claimRepository, &mockEmailRepository{}, administrator)

	// Then
	tx, err := blockchainUsecase.TransferCedarCoinsToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", "DUMMY_CLAIM_CODE")
	assert.NilError(t, err)
	assert.Equal(t, expectedTx.BlockNumber, tx.BlockNumber)
	assert.Equal(t, expectedTx.TxID, tx.TxID)
}

func TestBlockchainUsecase_TransferCedarCoinsToAddress_WithUnauthorizedClaim(t *testing.T) {
	// Given
	redeemedClaim := pkgmodel.NewClaim("Mike", "email@dummy.fr", 45)
	redeemedClaim.Status = pkgmodel.ClaimStatusClosed
	blockchainRepository := stubBlockchainRepository(nil, nil)
	claimRepository := stubClaimRepository(&redeemedClaim)
	administrator := pkgmodel.EmailRecipient{
		Address: "dummy@email.com",
		Name:    "Dummy Name",
	}
	blockchainUsecase := NewBlockchainUsecase(blockchainRepository, claimRepository, &mockEmailRepository{}, administrator)

	// Then
	_, err := blockchainUsecase.TransferCedarCoinsToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", "DUMMY_CLAIM_CODE")
	assert.ErrorContains(t, err, "not authorized to claim tokens")
}

func TestBlockchainUsecase_TransferCedarCoinsToAddress_WithBlockchainTransferError(t *testing.T) {
	// Given
	redeemedClaim := pkgmodel.NewClaim("Mike", "email@dummy.fr", 45)
	redeemedClaim.Status = pkgmodel.ClaimStatusOpen
	blockchainRepository := stubBlockchainRepository(nil, errors.New("blockchain transfer error"))
	claimRepository := stubClaimRepository(&redeemedClaim)
	administrator := pkgmodel.EmailRecipient{
		Address: "dummy@email.com",
		Name:    "Dummy Name",
	}
	blockchainUsecase := NewBlockchainUsecase(blockchainRepository, claimRepository, stubEmailRepository(), administrator)

	// Then
	_, err := blockchainUsecase.TransferCedarCoinsToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", "DUMMY_CLAIM_CODE")
	assert.ErrorContains(t, err, "blockchain transfer error")
}

func stubClaimRepository(claim *pkgmodel.Claim) *mockClaimRepository {
	mockClaimRepository := &mockClaimRepository{}
	mockClaimRepository.
		On("FindById", mock.Anything).
		Return(claim).
		Once()
	mockClaimRepository.
		On("UpdateById", mock.Anything).
		Return(nil)
	return mockClaimRepository
}

func stubBlockchainRepository(tx *model.Transaction, error error) *mockBlockchainRepository {
	mockBlockchainRepository := &mockBlockchainRepository{}
	mockBlockchainRepository.
		On("TransferCedarCoins", mock.Anything, mock.Anything).
		Return(tx, error)
	return mockBlockchainRepository
}

func stubEmailRepository() *mockEmailRepository {
	mockEmailRepository := &mockEmailRepository{}
	mockEmailRepository.
		On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
	return mockEmailRepository
}
