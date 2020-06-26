package rest

import (
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockClaimUsecase struct {
	mock.Mock
}

func (m *mockClaimUsecase) SaveClaim(claim *model.Claim) error {
	args := m.Called(claim)
	return args.Error(0)
}

func (m *mockClaimUsecase) FindClaims() ([]*model.Claim, error) {
	args := m.Called()
	return nil, args.Error(0)
}

func (m *mockClaimUsecase) FindClaim(id string) (*model.Claim, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Claim), nil
}

type mockHealthUsecase struct {
	mock.Mock
}

func (m *mockHealthUsecase) Health() *model.Health {
	args := m.Called()
	return args.Get(0).(*model.Health)
}

func withRouter(router *mux.Router, f func(t *testing.T, router *mux.Router)) func(t *testing.T) {
	return func(t *testing.T) {
		f(t, router)
	}
}

func executeRequest(req *http.Request, handler *mux.Router) *httptest.ResponseRecorder {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	return rr
}

func buildControllerWithClaimUseCaseError(error error, ucMethod string) *restController {
	mockHealthUsecase := &mockHealthUsecase{}
	mockClaimUsecase := &mockClaimUsecase{}
	mockClaimUsecase.
		On(ucMethod, mock.Anything).
		Return(error).
		Twice()
	rcc := NewRestController(mux.NewRouter(), mockClaimUsecase, mockHealthUsecase, "")
	rcc.Initialize()

	return rcc
}

func buildControllerWithClaimUseCaseReturning(claim *model.Claim, ucMethod string) *restController {
	mockHealthUsecase := &mockHealthUsecase{}
	mockClaimUsecase := &mockClaimUsecase{}
	mockClaimUsecase.
		On(ucMethod, mock.Anything).
		Return(claim).
		Twice()
	rcc := NewRestController(mux.NewRouter(), mockClaimUsecase, mockHealthUsecase, "")
	rcc.Initialize()

	return rcc
}

func buildControllerWitHealthUseCaseReturning(health *model.Health, ucMethod string) *restController {
	mockClaimUsecase := &mockClaimUsecase{}
	mockHealthUsecase := &mockHealthUsecase{}
	mockHealthUsecase.
		On(ucMethod, mock.Anything).
		Return(health).
		Twice()
	rcc := NewRestController(mux.NewRouter(), mockClaimUsecase, mockHealthUsecase, "")
	rcc.Initialize()

	return rcc
}

func getStringWithoutNewLine(toAssert string) string {
	return strings.TrimSuffix(toAssert, "\n")
}
