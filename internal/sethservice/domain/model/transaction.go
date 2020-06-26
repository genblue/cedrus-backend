package model

type Transaction struct {
	TxID        string `json:"txId" bson:"tx-id"`
	BlockNumber string `json:"blockNumber" bson:"block-number"`
}
