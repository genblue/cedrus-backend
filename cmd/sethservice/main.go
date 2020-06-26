package main

import (
	"context"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/infrastructure/blockchain"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/infrastructure/rest"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/usecases"
	pkgmodel "github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/genblue-private/cedrus-backend/pkg/infrastructure/database"
	"github.com/genblue-private/cedrus-backend/pkg/infrastructure/email"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
)

func main() {
	httpRouter := mux.NewRouter()
	mongoURI := getEnv("DB_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("DB_NAME", "")
	mongoClient := getMongoDBClient(mongoURI)

	// Repositories
	mongoClaimRepository := database.NewMongoClaimRepository(mongoClient, "claims", mongoDBName)
	tokenContractAddress := getEnv("TOKEN_CONTRACT_ADDRESS", "")
	log.Println("Contract's address is", tokenContractAddress)
	sethBlockchainRepository := blockchain.NewSethBlockchainRepository(tokenContractAddress)

	sendgridApiKey := getEnv("SENDGRID_API_KEY", "")
	if sendgridApiKey == "" {
		log.Fatal("ERROR", "could not get SENDGRID_API_KEY environment variable")
	}
	sendgridEmailRepository := email.NewSendgridEmailRepository(sendgridApiKey)

	// Use case
	administratorEmail := getEnv("ADMINISTRATOR_EMAIL", "")
	if administratorEmail == "" {
		log.Fatal("ERROR", "ADMINISTRATOR_EMAIL environment variable not set. Aborting.")
	}
	administratorName := getEnv("ADMINISTRATOR_NAME", "")
	if administratorName == "" {
		log.Fatal("ERROR", "ADMINISTRATOR_NAME environment variable not set. Aborting.")
	}
	administrator := pkgmodel.EmailRecipient{
		Address: administratorEmail,
		Name:    administratorName,
	}
	blockchainUsecase := usecases.NewBlockchainUsecase(sethBlockchainRepository, mongoClaimRepository, sendgridEmailRepository, administrator)

	// HTTP API
	startAPI(httpRouter, blockchainUsecase)
}

func startAPI(httpRouter *mux.Router, blockhainUsecase *usecases.BlockchainUsecase) {
	restController := rest.NewRestController(httpRouter, blockhainUsecase)
	restController.Initialize()
	port, err := strconv.Atoi(getEnv("API_PORT", "8002"))
	if err != nil {
		log.Fatal("ERROR", err)
	}
	restController.Run(port)
}

func getMongoDBClient(mongoURI string) *mongo.Client {
	// Set MongoDB client options and connect
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the MongoDB connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("ERROR", "could not connect to MongoDB")
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
	return client
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
