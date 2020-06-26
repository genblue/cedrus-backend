package rest

import (
	"encoding/json"
	"errors"
	inputs2 "github.com/genblue-private/cedrus-backend/internal/cedrusservice/infrastructure/rest/inputs"
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestRestController_PostClaim(t *testing.T) {
	// Given
	rcc := buildControllerWithClaimUseCaseError(nil, "SaveClaim")
	rccWithError := buildControllerWithClaimUseCaseError(errors.New("error from the DB"), "SaveClaim")

	jsonBody, _ := json.Marshal(inputs2.NewClaim{
		Name:      "Mike",
		Email:     "dummy@email.com",
		TreeCount: 666,
	})

	// Then
	t.Run("With expected JSON body", withRouter(rcc.router,
		testPostClaimFunc(string(jsonBody), http.StatusCreated, `{"data":null,"status":"success"}`)))
	t.Run("With empty body", withRouter(rcc.router,
		testPostClaimFunc("", http.StatusBadRequest, `{"data":"EOF","status":"fail"}`)))
	t.Run("With a use case error", withRouter(rccWithError.router,
		testPostClaimFunc(string(jsonBody), http.StatusBadRequest, `{"data":"error from the DB","status":"fail"}`)))
}

func testPostClaimFunc(jsonInput string, expectedStatus int, expectedResponse string) func(t *testing.T, router *mux.Router) {
	return func(t *testing.T, router *mux.Router) {
		request, err := http.NewRequest("POST", "/api/v1/claims", strings.NewReader(jsonInput))
		if err != nil {
			t.Fatal(err)
		}
		requestResponse := executeRequest(request, router)

		assert.Equal(t, expectedStatus, requestResponse.Code, "Bad status code")
		assert.Equal(t, expectedResponse, getStringWithoutNewLine(requestResponse.Body.String()), "Bad body")
	}
}

func TestRestController_GetClaims(t *testing.T) {
	// Given
	rcc := buildControllerWithClaimUseCaseError(nil, "FindClaims")
	rccWithError := buildControllerWithClaimUseCaseError(errors.New("error from the DB"), "FindClaims")

	// Then
	t.Run("Ok", withRouter(rcc.router,
		testGetClaimsFunc(http.StatusOK, `{"data":null,"status":"success"}`)))
	t.Run("With a use case error", withRouter(rccWithError.router,
		testGetClaimsFunc(http.StatusBadRequest, `{"data":"error from the DB","status":"fail"}`)))
}

func testGetClaimsFunc(expectedStatus int, expectedResponse string) func(t *testing.T, router *mux.Router) {
	return func(t *testing.T, router *mux.Router) {
		request, err := http.NewRequest("GET", "/api/v1/claims", nil)
		if err != nil {
			t.Fatal(err)
		}
		requestResponse := executeRequest(request, router)

		assert.Equal(t, expectedStatus, requestResponse.Code, "Bad status code")
		assert.Equal(t, expectedResponse, getStringWithoutNewLine(requestResponse.Body.String()), "Bad body")
	}
}

func TestRestController_GetClaim(t *testing.T) {
	// Given
	claim := model.NewClaim("Mike", "email@email.fr", 666)
	rcc := buildControllerWithClaimUseCaseReturning(&claim, "FindClaim")

	// When
	request, err := http.NewRequest("GET", "/api/v1/claims/"+claim.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	requestResponse := executeRequest(request, rcc.router)

	// Then
	jsonBody, err := json.Marshal(claim)
	assert.Equal(t, http.StatusOK, requestResponse.Code, "Bad status code")
	assert.Equal(t, `{"data":`+string(jsonBody)+`,"status":"success"}`, getStringWithoutNewLine(requestResponse.Body.String()), "Bad body")
}

func TestRestController_GetHealth(t *testing.T) {
	// Given
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	health := model.NewHealth(true)
	rcc := buildControllerWitHealthUseCaseReturning(&health, "Health")

	// When
	rr := executeRequest(req, rcc.router)

	// Then
	expectedResponse := `{"data":{"connectedToDb":true},"status":"success"}`
	assert.Equal(t, http.StatusOK, rr.Code, "Bad status code")
	assert.Equal(t, expectedResponse, getStringWithoutNewLine(rr.Body.String()), "Bad body")
}
