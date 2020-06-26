package model

import (
	"crypto/rand"
	"fmt"
	"time"
)

const (
	ClaimStatusNew       = 100 // claim status appears as New when it is created
	ClaimStatusOpen      = 200 // claim status appears as Open when email has been dispatched and it is awaiting settlement
	ClaimStatusApproved  = 300 // claim status is approved and will close shortly
	ClaimStatusRejected  = 400 // claim status changes to Rejected if the approvers reject the claim
	ClaimStatusClosed    = 500 // claim status appears as Closed after the claim is settled
	ClaimStatusCancelled = 600 // claim status can be updated to Cancelled only by internal staff
)

type Claim struct {
	ID              string `json:"_id" bson:"_id"`
	Name            string `json:"name" bson:"name"`
	Email           string `json:"email" bson:"email"`
	CreationDate    int64  `json:"creationDate" bson:"creation-date"`
	TreeCount       uint   `json:"treeCount" bson:"tree-count"`
	Status          int    `json:"status" bson:"status"`
	EmailSent       bool   `json:"emailSent" bson:"email-sent"`
	EmailSentDate   int64  `json:"emailSentDate" bson:"email-sent-date"`
	ClaimCode       string `json:"claimCode" bson:"claim-code"`
	SettlementDate  int64  `json:"settlementDate" bson:"settlement-date"`
	TransferAddress string `json:"transferAddress" bson:"transfer-address"`
	Memo            string `json:"memo" bson:"memo"`
}

// NewClaim  creates a claim object initialized with a generated token and timestamp
func NewClaim(name string, email string, treeCount uint) Claim {
	code := generateClaimCode()

	return Claim{
		Name:            name,
		Email:           email,
		Memo:            email,
		ID:              code,
		CreationDate:    time.Now().UnixNano(),
		TreeCount:       treeCount,
		ClaimCode:       code,
		Status:          ClaimStatusNew,
		EmailSent:       false,
		EmailSentDate:   0,
		SettlementDate:  0,
		TransferAddress: ""}
}

// generateClaimCode creates a random 8 byte token
func generateClaimCode() string {
	key := [8]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", key)
}
