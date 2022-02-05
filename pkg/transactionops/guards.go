package transactionops

import (
	"errors"
	"getherscan-cli/cmd/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GuardedSingleTransaction(conn *ethclient.Client, txHash common.Hash) error {
	tx, isPending, networkId := GetSingleTransaction(conn, txHash)

	if tx == nil || networkId == nil {
		message.OperationErrorMessage()
		return errors.New("GuardedSingleTransaction: Tx or NetworkId is nil")
	}

	ShowSimpleTransactionWithData(tx, isPending, networkId)
	return nil
}
