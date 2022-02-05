package blockops

import (
	"context"
	"getherscan-cli/cmd/message"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetSingleBlock(conn *ethclient.Client, blockNumber uint64) (*types.Block, *big.Int, error) {
	block, err := conn.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		message.NetworkFailWhileRequestMessage()
		return nil, nil, err
	}
	networkId, err := conn.NetworkID(context.Background())
	if err != nil {
		message.NetworkFailToGetNetworkId()
		return nil, nil, err
	}
	return block, networkId, nil
}
