package blockchain

import (
	"context"
	"errors"
	"fmt"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/model"
	"github.com/genblue-private/cedrus-backend/internal/sethservice/domain/repository"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type SethBlockchainRepository struct {
	tokenContractAddress string
}

func NewSethBlockchainRepository(tokenContractAddress string) repository.BlockchainRepository {
	return &SethBlockchainRepository{
		tokenContractAddress: tokenContractAddress,
	}
}

func (s SethBlockchainRepository) TransferCedarCoins(to string, amount uint) (transaction *model.Transaction, err error) {
	amountStr := strconv.FormatUint(uint64(amount), 10)
	cmdToExecute := fmt.Sprintf(
		"seth send %s \"mintCedarCoin(address,uint256)\" %s %s",
		s.tokenContractAddress,
		to,
		padRightSide(amountStr, 18))

	log.Printf("Executing transfer: %s", cmdToExecute)
	cmd := exec.CommandContext(context.Background(), "sh", "-c", cmdToExecute)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "invalid hexdata") ||
			strings.Contains(string(out), "bad address") {
			return nil, errors.New("invalid ethereum address")
		}
		if strings.Contains(string(out), "account not found") {
			return nil, errors.New("sender account for signing not found")
		}
		if strings.Contains(string(out), "insufficient funds") {
			return nil, errors.New("sender account has insufficient funds to send the transaction")
		}

		return nil, errors.New(string(out) + " (" + err.Error() + ")")
	}

	transaction = getTransactionDetails(out)

	return transaction, nil
}

func padRightSide(str string, count int) string {
	return str + strings.Repeat("0", count)
}

func getTransactionDetails(out []byte) *model.Transaction {
	lengthAsSizeOfString := -1
	txIDPattern := regexp.MustCompile("(0x.*?)\\n")
	txIDMatches := txIDPattern.FindAllStringSubmatch(string(out), lengthAsSizeOfString)

	blockNumberPattern := regexp.MustCompile("block(.*?).\\n")
	blockNumberMatches := blockNumberPattern.FindAllStringSubmatch(string(out), lengthAsSizeOfString)

	fullMatch := 0
	matchGroup := 1
	transaction := &model.Transaction{
		TxID:        txIDMatches[fullMatch][matchGroup],
		BlockNumber: blockNumberMatches[fullMatch][matchGroup],
	}
	return transaction
}

func (s SethBlockchainRepository) Ping() error {
	cmd := exec.CommandContext(context.Background(), "sh", "-c", "seth")
	_, err := cmd.CombinedOutput()
	log.Println("pinged seth, all's fine")
	if err != nil {
		return err
	}

	return nil
}

func (s SethBlockchainRepository) AccountBalance() (*model.AccountBalance, error) {
	cmd := exec.CommandContext(context.Background(), "sh", "-c", "seth ls")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	accountBalance, err := getAccountBalance(out)
	if err != nil {
		return nil, err
	}

	log.Println("account balance", accountBalance)
	return accountBalance, nil
}

func getAccountBalance(out []byte) (*model.AccountBalance, error) {
	lengthAsSizeOfString := -1

	accountPattern := regexp.MustCompile("^(.*?)\\t")
	accountMatches := accountPattern.FindAllStringSubmatch(string(out), lengthAsSizeOfString)
	if accountMatches == nil {
		return nil, errors.New("no ethereum account defined")
	}
	balancePattern := regexp.MustCompile("\\t([[:alnum:]])\\n")
	balanceMatches := balancePattern.FindAllStringSubmatch(string(out), lengthAsSizeOfString)

	fullMatch := 0
	matchGroup := 1
	accountBalance := model.AccountBalance{
		Address: accountMatches[fullMatch][matchGroup],
		Balance: balanceMatches[fullMatch][matchGroup],
	}

	return &accountBalance, nil
}
