package usecases

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestEmailUsecase_SendEmailsToNewClaims(t *testing.T) {
	// Given
	claims := getDummyClaims()
	mockClaimRepository := stubClaimRepository(claims)
	mockEmailRepository := stubEmailRepository()
	emailUsecase := NewEmailUsecase(mockClaimRepository, mockEmailRepository)

	// Then
	err := emailUsecase.SendEmailsToNewClaims()
	if err != nil {
		t.Error(err)
	}
}

func stubClaimRepository(claims []*model.Claim) *mockClaimRepository {
	mockClaimRepository := &mockClaimRepository{}
	mockClaimRepository.
		On("FindAllByEmailUnsent", mock.Anything).
		Return(claims).
		Once()
	mockClaimRepository.
		On("UpdateById", mock.Anything).
		Return(nil).
		Times(2)
	return mockClaimRepository
}

func stubEmailRepository() *mockEmailRepository {
	mockEmailRepository := &mockEmailRepository{}
	mockEmailRepository.
		On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Times(2)
	return mockEmailRepository
}

func getDummyClaims() []*model.Claim {
	firstClaim := model.NewClaim("Mike", "email@dummy.fr", 45)
	secondClaim := model.NewClaim("Roger", "email45@dummy.fr", 666)
	var claims []*model.Claim
	return append(claims, &firstClaim, &secondClaim)
}
