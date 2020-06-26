package repository

import "github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"

type BlockchainRepository interface {
	TransferCedarCoins(to string, amount uint) (transaction *model.Transaction, err error)
	AccountBalance() (*model.AccountBalance, error)
	Ping() error
}
