package main

import (
	"context"
	"github.com/genblue-private/cedrus-backend/internal/emailservice/usecases"
	"github.com/genblue-private/cedrus-backend/pkg/infrastructure/database"
	"github.com/genblue-private/cedrus-backend/pkg/infrastructure/email"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	mongoURI := getEnv("DB_URI", "mongodb://localhost:27017")
	mongoDBName := getEnv("DB_NAME", "")
	mongoClient := getMongoDBClient(mongoURI)

	// Repositories
	mongoClaimRepository := database.NewMongoClaimRepository(mongoClient, "claims", mongoDBName)
	sendgridApiKey := getEnv("SENDGRID_API_KEY", "")
	if sendgridApiKey == "" {
		log.Fatal("ERROR", "could not get SENDGRID_API_KEY environment variable")
	}
	sendgridEmailRepository := email.NewSendgridEmailRepository(sendgridApiKey)

	// Use case
	emailUsecase := usecases.NewEmailUsecase(mongoClaimRepository, sendgridEmailRepository)
	err := emailUsecase.SendEmailsToNewClaims()
	if err != nil {
		log.Fatal(err)
	}
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
