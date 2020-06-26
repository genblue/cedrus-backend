package usecases

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"gotest.tools/assert"
	"testing"
)

func TestClaimUsecase_SaveClaim_WithBadEmailFormat(t *testing.T) {
	// Given
	claim := model.NewClaim("Mike", "Bad email", 45)
	claimUsecase := NewClaimUsecase(&mockClaimRepository{})

	// When
	err := claimUsecase.SaveClaim(&claim)

	// Then
	assert.Error(t, err, "invalid format: Bad email")
}

/* Email host is commented due to instability
func TestClaimUsecase_SaveClaim_WithBadEmailHost(t *testing.T) {
    // Given
    claim := model.NewClaim("Mike", "fake@dummyhost.verydummy", 45)
    claimUsecase := NewClaimUsecase(&mockClaimRepository{})

    // When
    err := claimUsecase.SaveClaim(&claim)

    // Then
    assert.Error(t, err, "unresolvable host: fake@dummyhost.verydummy")
}
*/
