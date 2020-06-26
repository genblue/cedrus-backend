package main

import (
	"context"
	"github.com/genblue-private/cedrus-backend/internal/cedrusservice/infrastructure/rest"
	"github.com/genblue-private/cedrus-backend/internal/cedrusservice/usecases"
	"github.com/genblue-private/cedrus-backend/pkg/infrastructure/database"
	pkgusecases "github.com/genblue-private/cedrus-backend/pkg/usecases"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
)

func main() {
	mongoURI := getEnv("DB_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("DB_NAME", "")
	httpRouter := mux.NewRouter()
	mongoClient := getMongoDBClient(mongoURI)

	// Repositories
	mongoClaimRepository := database.NewMongoClaimRepository(mongoClient, "claims", mongoDBName)

	// Use cases
	claimUsecase := usecases.NewClaimUsecase(mongoClaimRepository)
	healthUsecase := pkgusecases.NewHealthUsecase(mongoClaimRepository)

	// HTTP API
	startAPI(httpRouter, claimUsecase, healthUsecase)
}

func startAPI(httpRouter *mux.Router, claimUsecase *usecases.ClaimUsecase, healthUsecase *pkgusecases.HealthUsecase) {
	sethServiceURL := getEnv("SETH_SERVICE_URL", "")
	log.Println("Seth service URL is", sethServiceURL)
	restController := rest.NewRestController(httpRouter, claimUsecase, healthUsecase, sethServiceURL)
	restController.Initialize()
	port, err := strconv.Atoi(getEnv("API_PORT", "8000"))
	if err != nil {
		log.Fatal("ERROR", err)
	}
	restController.Run(port)
}

func getMongoDBClient(mongoURI string) *mongo.Client {
	// Set MongoDB client options and connect
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the MongoDB connection
	err = client.Ping(context.TODO(), nil)
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
