package database

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMongoClaimRepository_FindById(t *testing.T) {
	// Given
	pool := getDockerPool()
	mongoClient, resource := GetMongoDBClient(t, *pool)
	defer pool.Purge(resource)

	mcr := NewMongoClaimRepository(mongoClient, "claims", "cedrus")
	expectedClaim := model.NewClaim("Mike", "email@email.fr", 2)
	err := mcr.Save(&expectedClaim)
	if err != nil {
		t.Error(err)
	}

	// When
	claim, err := mcr.FindById(expectedClaim.ID)
	if err != nil {
		t.Error(err)
	}

	// Then
	assert.Equal(t, &expectedClaim, claim, "Bad claim")
}

func TestMongoClaimRepository_FindAll(t *testing.T) {
	// Given
	pool := getDockerPool()
	mongoClient, resource := GetMongoDBClient(t, *pool)
	defer pool.Purge(resource)

	mcr := NewMongoClaimRepository(mongoClient, "claims", "cedrus")
	expectedFirstClaim := model.NewClaim("Mike", "email@dummy.com", 2)
	expectedSecondClaim := model.NewClaim("Roger", "emailTwo@dummy.com", 9)
	_ = mcr.Save(&expectedFirstClaim)
	_ = mcr.Save(&expectedSecondClaim)

	// When
	claims, err := mcr.FindAll()
	if err != nil {
		t.Error(err)
	}

	// Then
	firstClaim, secondClaim := claims[0], claims[1]
	assert.Equal(t, &expectedFirstClaim, firstClaim, "Bad claim")
	assert.Equal(t, &expectedSecondClaim, secondClaim, "Bad claim")
}

func TestMongoClaimRepository_FindAllByEmailUnsent(t *testing.T) {
	// Given
	pool := getDockerPool()
	mongoClient, resource := GetMongoDBClient(t, *pool)
	defer pool.Purge(resource)

	mcr := NewMongoClaimRepository(mongoClient, "claims", "cedrus")
	claimWithUnsentEmail := model.NewClaim("Mike", "email@dummy.com", 2)
	_ = mcr.Save(&claimWithUnsentEmail)
	claimWithSentEmail := model.NewClaim("Roger", "emailTwo@dummy.com", 9)
	claimWithSentEmail.EmailSent = true
	_ = mcr.Save(&claimWithSentEmail)

	// When
	claims, err := mcr.FindAllByEmailUnsent()
	if err != nil {
		t.Error(err)
	}

	// Then
	assert.Len(t, claims, 1)
	assert.Equal(t, &claimWithUnsentEmail, claims[0], "Bad claim")
}

func TestMongoClaimRepository_UpdateById(t *testing.T) {
	// Given
	pool := getDockerPool()
	mongoClient, resource := GetMongoDBClient(t, *pool)
	defer pool.Purge(resource)

	mcr := NewMongoClaimRepository(mongoClient, "claims", "cedrus")
	claimWithUnsentEmail := model.NewClaim("Mike", "email@dummy.com", 2)
	_ = mcr.Save(&claimWithUnsentEmail)

	// When
	claimWithUnsentEmail.EmailSent = true
	claimWithUnsentEmail.EmailSentDate = time.Now().Unix()
	err := mcr.UpdateById(&claimWithUnsentEmail)
	if err != nil {
		t.Error(err)
	}

	// Then
	claim, err := mcr.FindById(claimWithUnsentEmail.ID)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, &claimWithUnsentEmail, claim, "Bad claim")
}
