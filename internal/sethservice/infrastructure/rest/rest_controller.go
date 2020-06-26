package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gamegos/jsend"
	_ "github.com/genblue-private/cedrus-backend/api-docs/sethservice"
	"github.com/genblue-private/cedrus-backend/internal/cedrusservice/infrastructure/rest/inputs"
	_ "github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/usecases"
	_ "github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

type restController struct {
	router *mux.Router
	buc    usecases.BlockchainUsecaseInterface
}

func NewRestController(
	router *mux.Router,
	buc usecases.BlockchainUsecaseInterface) *restController {
	return &restController{
		router: router.PathPrefix("/api/").Subrouter(),
		buc:    buc,
	}
}

// @title Seth service API
// @version v1
// @description For managing blockchain transfers

// @contact.name API Support
// @contact.email email@ded.fr

// @BasePath /api/v1
func (rcc *restController) Initialize() {
	rcc.router.HandleFunc("/v1/health", rcc.GetHealth).Methods("GET")
	rcc.router.HandleFunc("/v1/transfer", rcc.PostTransfer).Methods("POST")
	rcc.router.HandleFunc("/v1/accounts_balance", rcc.GetAccountsBalance).Methods("GET")
	rcc.router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler).Methods("GET")

	_ = rcc.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		fmt.Println(methods, path)
		return nil
	})
}

func (rcc *restController) Run(port int) {
	log.Printf("Running HTTP API on port %v...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), rcc.router))
}

// GetHealth godoc
// @Summary Get health of the service
// @Success 200
// @Router /health [get]
func (rcc *restController) GetHealth(w http.ResponseWriter, r *http.Request) {
	err := rcc.buc.Health()
	if err != nil {
		_, _ = jsend.Wrap(w).Data(err).Status(http.StatusInternalServerError).Send()
	}

	_, _ = jsend.Wrap(w).Status(http.StatusOK).Send()
}

// GetHealth godoc
// @Summary Get balance of the ethereum accounts
// @Success 200 {object} model.AccountBalance
// @Router /accounts_balance [get]
func (rcc *restController) GetAccountsBalance(w http.ResponseWriter, r *http.Request) {
	accountBalance, err := rcc.buc.FindAccountBalance()
	if err != nil {
		_, _ = jsend.Wrap(w).Data(err).Status(http.StatusInternalServerError).Send()
	}

	_, _ = jsend.Wrap(w).Data(accountBalance).Status(http.StatusOK).Send()
}

// PostTransfer godoc
// @Summary Transfer Cedar coins
// @Param transfer body inputs.NewTransfer true "Transfer"
// @Success 202
// @Router /transfer [post]
func (rcc *restController) PostTransfer(w http.ResponseWriter, r *http.Request) {
	log.Println("Received", r.Body, "to transfer")
	var inputTransfer inputs.NewTransfer
	errDecoding := json.NewDecoder(r.Body).Decode(&inputTransfer)
	if errDecoding != nil {
		log.Println("ERROR", errDecoding)
		_, _ = jsend.Wrap(w).Data(errDecoding.Error()).Status(http.StatusBadRequest).Send()
		return
	}

	txID, err := rcc.buc.TransferCedarCoinsToAddress(inputTransfer.Address, inputTransfer.ClaimCode)
	if err != nil {
		log.Println("ERROR", err)
		_, _ = jsend.Wrap(w).Data(err.Error()).Status(http.StatusBadRequest).Send()
		return
	}

	_, _ = jsend.Wrap(w).Data(txID).Status(http.StatusAccepted).Send()
}
