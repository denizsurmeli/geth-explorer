package transactionops

import (
	"fmt"
	"getherscan-cli/cmd/message"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
)

func ShowSimpleTransaction(transaction *types.Transaction, networkId *big.Int) {
	asMessage, err := transaction.AsMessage(types.NewLondonSigner(networkId), nil)
	if err != nil {
		message.OperationErrorMessage()
	}
	fromToString := asMessage.From().String()
	var toToString string
	if asMessage.To() == nil {
		toToString = "[?]"
	} else {
		toToString = asMessage.To().String()
	}
	valueToEtherThenString := asMessage.Value().String()
	gasPaid := asMessage.Gas() * asMessage.GasPrice().Uint64()
	totalValue := transaction.Cost().Uint64()
	fmt.Printf("[TX] %s ----> %s \n\t| Value:%s wei \n\t| Gas Paid:%d wei \n\t| Total Value:%d wei\n",
		fromToString,
		toToString,
		valueToEtherThenString,
		gasPaid,
		totalValue)
}

func ShowSimpleTransactionWithData(tx *types.Transaction, isPending bool, networkId *big.Int) {
	txAsMessage, err := tx.AsMessage(types.NewLondonSigner(networkId), nil)
	if err != nil {
		message.OperationErrorMessage()
	}
	color.Cyan("[TRANSACTION] Tx Hash: %s Is Pending ?:%t", tx.Hash().String(), isPending)
	color.Yellow("\t[FROM]:%s --> [TO]:%s", txAsMessage.From().String(), txAsMessage.To().String())
	color.Blue("\t[VALUE] Gas Limit:%d | Value:%d ", txAsMessage.Gas(), txAsMessage.Value().Uint64())
	fmt.Println("\t[DATA] =", "0x"+hexutils.BytesToHex(txAsMessage.Data()))
}
