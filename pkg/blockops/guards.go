package blockops

import (
	"errors"
	"getherscan-cli/cmd/message"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GuardedSingleBlock(conn *ethclient.Client, blocknumber uint64) error {
	err := GuardedSingleBlockWithMetadata(conn, blocknumber)
	if err != nil {
		message.OperationErrorMessage()
		return errors.New("GuardedSingleBlock:Metadata error")
	}

	err = GuardedSingleBlockWithTransactions(conn, blocknumber)
	if err != nil {
		message.OperationErrorMessage()
		return errors.New("GuardedSingleBlock:Transactions error")
	}
	return nil

}
func GuardedSingleBlockWithMetadata(conn *ethclient.Client, blocknumber uint64) error {
	block, networkId, err := GetSingleBlock(conn, blocknumber)
	if err != nil {
		message.OperationErrorMessage()
		return errors.New("GuardedSingleBlock: err is not nil")
	}
	ShowSingleBlockMetadata(block, networkId)
	return nil

}

func GuardedSingleBlockWithTransactions(conn *ethclient.Client, blocknumber uint64) error {
	block, _, err := GetSingleBlock(conn, blocknumber)
	if err != nil {
		message.OperationErrorMessage()
		return errors.New("GuardedSingleBlock: err is not nil")
	}
	ShowSingleBlockWithSimpleTransactions(conn, block)
	return nil

}
