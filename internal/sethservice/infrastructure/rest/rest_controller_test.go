package rest

import (
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestRestController_PostTransfer(t *testing.T) {
	// Given
	jsonInput := "{\"address\": \"0x0808334180392c61b15065ad56130b3b35a22806\", \"claimCode\": \"1e00167939cb6694\"}"
	req, err := http.NewRequest("POST", "/api/v1/transfer", strings.NewReader(jsonInput))
	if err != nil {
		t.Fatal(err)
	}
	expectedTx := model.Transaction{
		TxID:        "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
		BlockNumber: "1234",
	}
	rcc := buildControllerWithBlockchainUseCaseReturning(&expectedTx, "TransferCedarCoinsToAddress")

	// When
	rr := executeRequest(req, rcc.router)

	// Then
	expectedResponse := `{"data":{"txId":"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae","blockNumber":"1234"},"status":"success"}`
	assert.Equal(t, http.StatusAccepted, rr.Code, "Bad status code")
	assert.Equal(t, expectedResponse, getStringWithoutNewLine(rr.Body.String()), "Bad body")
}
