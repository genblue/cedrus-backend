package model

type AccountBalance struct {
	Address string `json:"address" bson:"address"`
	Balance string `json:"balance" bson:"balance"`
}
