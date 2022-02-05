package transactionops

import (
	"context"
	"getherscan-cli/cmd/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetSingleTransaction(conn *ethclient.Client, txHash common.Hash) (*types.Transaction, bool, *big.Int) {
	tx, isPending, err := conn.TransactionByHash(context.Background(), txHash)
	if tx == nil || err != nil {
		message.OperationFailTransaction()
		return nil, false, nil
	}

	networkId, err := conn.NetworkID(context.Background())
	if err != nil {
		message.NetworkFailToGetNetworkId()
		return tx, false, nil
	}
	return tx, isPending, networkId
}
