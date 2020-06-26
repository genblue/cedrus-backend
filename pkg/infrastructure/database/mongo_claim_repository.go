package database

import (
	"context"
	"github.com/genblue-private/cedrus-backend/pkg/domain/model"
	"github.com/genblue-private/cedrus-backend/pkg/domain/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type mongoClaimRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	context    context.Context
}

func NewMongoClaimRepository(client *mongo.Client, collectionName string, databaseName string) repository.ClaimRepository {
	return &mongoClaimRepository{
		client:     client,
		collection: client.Database(databaseName).Collection(collectionName),
		context:    context.Background(),
	}
}

func (m mongoClaimRepository) FindAll() ([]*model.Claim, error) {
	var claims []*model.Claim

	cur, err := m.collection.Find(m.context, bson.M{})
	if err != nil {
		return nil, err
	}

	for cur.Next(m.context) {
		var claim model.Claim
		err = cur.Decode(&claim)
		if err != nil {
			log.Printf("error while decoding the document: %s", err)
		} else {
			claims = append(claims, &claim)
		}
	}
	return claims, nil
}

func (m mongoClaimRepository) FindById(id string) (*model.Claim, error) {
	var claim *model.Claim

	idFilter := bson.D{{"_id", id}}
	err := m.collection.FindOne(m.context, idFilter).Decode(&claim)
	if err != nil {
		return nil, err
	}

	return claim, nil
}

func (m mongoClaimRepository) Save(claim *model.Claim) error {
	_, err := m.collection.InsertOne(m.context, claim)
	if err != nil {
		return err
	}

	return nil
}

func (m mongoClaimRepository) UpdateById(claim *model.Claim) error {
	idFilter := bson.M{
		"_id": bson.M{
			"$eq": claim.ID,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"email-sent":       claim.EmailSent,
			"email-sent-date":  claim.EmailSentDate,
			"status":           claim.Status,
			"settlement-date":  claim.SettlementDate,
			"transfer-address": claim.TransferAddress,
		},
	}

	res, err := m.collection.UpdateOne(m.context, idFilter, update)
	if err != nil {
		return err
	}

	logrus.Debug("Updated ", res.UpsertedCount, " row with ", claim)
	return nil
}

func (m mongoClaimRepository) Ping() error {
	err := m.client.Ping(m.context, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

func (m mongoClaimRepository) FindAllByEmailUnsent() ([]*model.Claim, error) {
	var claims []*model.Claim

	emailNotSentFilter := bson.D{{"email-sent", false}}
	cur, err := m.collection.Find(m.context, emailNotSentFilter)
	if err != nil {
		return nil, err
	}

	for cur.Next(m.context) {
		var claim model.Claim
		err = cur.Decode(&claim)
		if err != nil {
			log.Printf("error while decoding the document: %s", err)
		} else {
			claims = append(claims, &claim)
		}
	}
	return claims, nil
}
