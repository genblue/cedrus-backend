package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gamegos/jsend"
	_ "github.com/genblue-private/cedrus-backend/api-docs/cedrusservice"
	"github.com/genblue-private/cedrus-backend/internal/cedrusservice/infrastructure/rest/inputs"
	"github.com/genblue-private/cedrus-backend/internal/cedrusservice/usecases"
	_ "github.com/genblue-private/cedrus-backend/internal/sethservice/infrastructure/rest/inputs"
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	pkgusecases "github.com/genblue-private/cedrus-backend/pkg/usecases"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type restController struct {
	router         *mux.Router
	cuc            usecases.ClaimUsecaseInterface
	huc            pkgusecases.HealthUsecaseInterface
	httpClient     *http.Client
	sethServiceURL string
}

func NewRestController(
	router *mux.Router,
	cuc usecases.ClaimUsecaseInterface,
	huc pkgusecases.HealthUsecaseInterface,
	sethServiceURL string) *restController {
	return &restController{
		router: router,
		cuc:    cuc,
		huc:    huc,
		httpClient: &http.Client{
			Timeout: 600 * time.Second,
		},
		sethServiceURL: sethServiceURL,
	}
}

// @title Cedrus service API
// @version v1
// @description For managing claims

// @contact.name API Support
// @contact.email email@ded.fr

// @BasePath /api/v1
func (rcc *restController) Initialize() {
	indexServer := http.FileServer(http.Dir("./static"))
	rcc.router.Handle("/", indexServer)
	rcc.router.Handle("/index.html", indexServer)
	rcc.router.Handle("/claim.html", indexServer)
	rcc.router.Handle("/spend.html", indexServer)
	rcc.router.Handle("/checkout.html", indexServer)
	rcc.router.Handle("/checkout-with-dai.html", indexServer)

	assetsServer := http.FileServer(http.Dir("./static/assets"))
	rcc.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsServer))

	rcc.router.HandleFunc("/api/v1/health", rcc.GetHealth).Methods("GET")
	rcc.router.HandleFunc("/api/v1/claims/{id}", rcc.GetClaim).Methods("GET")
	rcc.router.HandleFunc("/api/v1/claims", rcc.GetClaims).Methods("GET")
	rcc.router.HandleFunc("/api/v1/claims", rcc.PostClaim).Methods("POST")
	rcc.router.HandleFunc("/api/v1/transfer", rcc.PostTransfer).Methods("POST")
	rcc.router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler).Methods("GET")

	_ = rcc.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		log.Println(methods, path)
		return nil
	})
}

func (rcc *restController) Run(port int) {
	log.Printf("Running HTTP API on port %v...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), rcc.router))
}

// GetHealth godoc
// @Summary get application health
// @Success 200 {object} model.Health
// @Router /health [get]
func (rcc *restController) GetHealth(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	health := rcc.huc.Health()
	_, _ = jsend.Wrap(w).Data(health).Status(http.StatusOK).Send()
}

// PostClaim godoc
// @Summary Create a Claim
// @Produce json
// @Param new-claim body inputs.NewClaim true "New Claim"
// @Success 201
// @Failure 400
// @Router /claims [post]
func (rcc *restController) PostClaim(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var inputClaim inputs.NewClaim
	errDecoding := json.NewDecoder(r.Body).Decode(&inputClaim)
	if errDecoding != nil {
		_, _ = jsend.Wrap(w).Data(errDecoding.Error()).Status(http.StatusBadRequest).Send()
		return
	}

	claim := model.NewClaim(inputClaim.Name, inputClaim.Email, inputClaim.TreeCount)
	errSaving := rcc.cuc.SaveClaim(&claim)
	if errSaving != nil {
		_, _ = jsend.Wrap(w).Data(errSaving.Error()).Status(http.StatusBadRequest).Send()
		return
	}

	_, _ = jsend.Wrap(w).Data(nil).Status(http.StatusCreated).Send()
}

// GetClaim godoc
// @Summary Get a Claim
// @Produce json
// @Param id path string true "Claim ID"
// @Success 200 {object} model.Claim
// @Router /claims/{id} [get]
func (rcc *restController) GetClaim(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	vars := mux.Vars(r)
	id := vars["id"]

	claim, err := rcc.cuc.FindClaim(id)
	if err != nil {
		_, _ = jsend.Wrap(w).Data(err.Error()).Status(http.StatusBadRequest).Send()
		return
	}

	_, _ = jsend.Wrap(w).Data(claim).Status(http.StatusOK).Send()
}

// GetClaims godoc
// @Summary Get all Claims
// @Produce json
// @Success 200 {object} model.Claim[]
// @Router /claims [get]
func (rcc *restController) GetClaims(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	claims, err := rcc.cuc.FindClaims()
	if err != nil {
		_, _ = jsend.Wrap(w).Data(err.Error()).Status(http.StatusBadRequest).Send()
		return
	}

	_, _ = jsend.Wrap(w).Data(claims).Status(http.StatusOK).Send()
}

// PostTransfer godoc
// @Summary Transfer Cedar coins
// @Param transfer body inputs.NewTransfer true "Transfer"
// @Success 202
// @Router /transfer [post]
func (rcc *restController) PostTransfer(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	response, err := rcc.httpClient.Post(rcc.sethServiceURL, "application/json", r.Body)
	if err != nil {
		log.Println("ERROR", "Could not proxy transfer to sethservice:", err)
		_, _ = jsend.Wrap(w).Data(err.Error()).Status(http.StatusInternalServerError).Send()
		return
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("ERROR", "could not read response body", err)
	}

	log.Printf("Successfully proxyied transfer to sethservice. Response: %s", json.RawMessage(responseBody))
	w.WriteHeader(response.StatusCode)
	_, _ = io.WriteString(w, string(json.RawMessage(responseBody)))
}

func logRequest(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println("ERROR", err)
	}
	log.Println(string(requestDump))
}
