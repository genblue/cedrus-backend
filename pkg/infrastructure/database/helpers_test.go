package database

import (
	"context"
	"fmt"
	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
)

func GetMongoDBClient(t *testing.T, pool dockertest.Pool) (*mongo.Client, *dockertest.Resource) {
	ctx := context.Background()
	var db *mgo.Session
	var err error

	resource, err := pool.Run("mongo", "3.4.23", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = mgo.Dial(fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Error(err)
	}

	// Check the MongoDB connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}

	return client, resource
}

func getDockerPool() *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool
}
