package rest

import (
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockBlockchainUsecase struct {
	mock.Mock
}

func (m *mockBlockchainUsecase) TransferCedarCoinsToAddress(to string, claimCode string) (transaction *model.Transaction, err error) {
	args := m.Called(to, claimCode)
	return args.Get(0).(*model.Transaction), nil
}

func (m *mockBlockchainUsecase) Health() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockBlockchainUsecase) FindAccountBalance() (*model.AccountBalance, error) {
	args := m.Called()
	return args.Get(0).(*model.AccountBalance), nil
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

func buildControllerWithBlockchainUseCaseReturning(tx *model.Transaction, ucMethod string) *restController {
	mockBlockchainUsecase := &mockBlockchainUsecase{}
	mockBlockchainUsecase.
		On(ucMethod, mock.Anything, mock.Anything).
		Return(tx).
		Twice()
	rcc := NewRestController(mux.NewRouter(), mockBlockchainUsecase)
	rcc.Initialize()

	return rcc
}

func getStringWithoutNewLine(toAssert string) string {
	return strings.TrimSuffix(toAssert, "\n")
}
