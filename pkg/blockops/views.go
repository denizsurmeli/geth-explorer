package blockops

import (
	"context"
	"fmt"
	"getherscan-cli/pkg/transactionops"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"math/big"
)

func ShowSingleBlockMetadata(block *types.Block, networkId *big.Int) {
	color.Green("[BLOCK] Block %d:\n", block.Header().Number.Uint64())
	color.Magenta("\t[BLOCK] Block hash:%s\n", block.Header().Hash().String())
	color.Yellow("\t[BLOCK] Block size:%s\n", block.Header().Size().String())
	color.Yellow("\t[BLOCK] Tx Count:%d\n", len(block.Transactions()))

}

func ShowSingleBlockWithSimpleTransactions(conn *ethclient.Client, block *types.Block) {
	color.Cyan("\t[BLOCK] Transactions(Simple View):\n")
	for _, transaction := range block.Transactions() {
		networkId, err := conn.NetworkID(context.Background())
		if err != nil {
			color.Red("[OPERATION] Chain ID could not be fetched while formatting the block info.Halting the exection.")
			fmt.Println(err.Error())
		}
		transactionops.ShowSimpleTransaction(transaction, networkId)
	}
}
