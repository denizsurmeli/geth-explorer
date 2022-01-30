package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/status-im/keycard-go/hexutils"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func main() {
	networkNamePtr := flag.String("network", "mainnet", "a string")
	operationPtr := flag.String("operation", "listen_headers", "a string")
	txHashPtr := flag.String("txhash", "0x0", "a string")

	// only use websockets
	flag.Parse()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		color.Red("[SETUP] Error while reading the .env file. Please check your .env file for correct setup.")
		panic("")
	}
	value, ok := viper.Get("INFURA_KEY").(string)
	if !ok {
		color.Red("[SETUP] Error while reading the .env file. Please check your .env file for correct setup.")
		panic("")
	}

	networkConfig := []string{"wss://", *networkNamePtr, ".infura.io/ws/v3/", value}
	color.Green("[NETWORK] Dialing Ethereum(%s) Node...", *networkNamePtr)

	conn, err := ethclient.Dial(strings.Join(networkConfig, ""))
	if err != nil {
		color.Red("[NETWORK] Error! Could not dial the node. Halting the execution.")
		panic("could not dial the node.")
	}
	// relative connection test, if we have dialed the wrong endpoint,
	// any function of eth.Client should behave unexpectedly.
	networkId, err := conn.NetworkID(context.Background())
	if networkId == nil || err != nil {
		color.Red("[NETWORK] Error! Could not dial the node. Halting the execution.")
		panic("could not dial the node.")
	}
	color.Green("[NETWORK] Node connection established.")
	switch *operationPtr {
	case "listen_headers":
		ListenHeaders(conn)
	case "listen_blocks":
		ListenBlocks(conn)
	case "lens_transaction":
		if *txHashPtr == "0x0" {
			color.Red("[PARAMETERERR] Tx Hash invalid or not given for this operation.Pass it by --txhash=<transaction hash>.")
			panic("")
		}
		LensTransaction(conn, *txHashPtr)

	}
}

func LensTransaction(conn *ethclient.Client, transaction string) {
	txHash := common.HexToHash(transaction)
	tx, isPending, err := conn.TransactionByHash(context.Background(), txHash)
	networkId, err := conn.NetworkID(context.Background())
	if err != nil {
		color.Red("[OPERATION] Network Id could not be fetched.Halting the execution")
		panic("")
	}

	txAsMessage, err := tx.AsMessage(types.NewLondonSigner(networkId), nil)
	if err != nil {
		color.Red("[OPERATION] Formatting as message failed. Halting the execution.")
		panic("")
	}
	if err != nil {
		color.Red("[OPERATION] Error occured while fetching the transaction. Maybe check the transaction hash ?")
		panic("")
	}
	color.Cyan("[TRANSACTION] Tx Hash: %s Is Pending ?:%t", tx.Hash().String(), isPending)
	color.Green("[FROM]:%s --> [TO]:%s", txAsMessage.From().String(), txAsMessage.To().String())
	color.Blue("Gas Limit:%d | Value:%d ", txAsMessage.Gas(), txAsMessage.Value().Uint64())
	color.Yellow("[TXDATA(bytes32)] Transaction data:")
	fmt.Println("\t[TXDATA] =", txAsMessage.Data())
	color.Green("\tFormat to Hex ->")
	fmt.Println("\t[HEXTXDATA] =", "0x"+hexutils.BytesToHex(txAsMessage.Data()))
}
func ListenBlocks(conn *ethclient.Client) {
	headers := make(chan *types.Header)
	subscription, err := conn.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		color.Red("[OPERATION] Listening to blocks behaved unexpectedly.Halting the execution.")
		fmt.Println(err.Error())
	}
	for {
		select {
		case err := <-subscription.Err():
			log.Fatal(err)
			panic("Something went wrong in the channel")
		case blockHeader := <-headers:
			t := time.Now()
			block, err := conn.BlockByHash(context.Background(), blockHeader.Hash())
			if err != nil {
				color.Red("[OPERATION] Block could not be fetched.Halting the execution.")
				fmt.Println(err.Error())
			}
			color.Green("[BLOCK] %d-%02d-%02dT%02d:%02d:%02d-00:00 Block %d found.\n", t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second(), block.Header().Number.Uint64())
			color.Cyan("[SIZE] Block size:%s\n", block.Header().Size().String())
			color.Magenta("[BLOCKINFO] Block hash:%s\n", block.Header().Hash().String())
			color.Blue("[BLOCKINFO] Tx Count:%d\n", len(block.Transactions()))
			color.Cyan("[SIMPLETXLIST] Transactions:\n")
			for index, transaction := range block.Transactions() {
				chainId, err := conn.NetworkID(context.Background())
				if err != nil {
					color.Red("[OPERATION] Chain ID could not be fetched while formatting the block info.Halting the exection.")
					fmt.Println(err.Error())
				}
				message, err := transaction.AsMessage(types.NewLondonSigner(chainId), nil)
				if err != nil {
					color.Red("[OPERATION] Message could not be fetched while formatting the block info.Halting the exection.")
					fmt.Println(err.Error())
				}
				fromToString := message.From().String()
				var toToString string
				if message.To() == nil {
					toToString = "?(Probably Signer Error.)"
				} else {
					toToString = message.To().String()
				}
				valueToEtherThenString := message.Value().String()
				gasPaid := (message.Gas() * message.GasPrice().Uint64())
				totalValue := transaction.Cost().Uint64()
				fmt.Printf("%d) %s ----> %s \n\t| Value:%s wei \n\t| Gas Paid:%d wei \n\t| Total Value:%d wei\n",
					index+1,
					fromToString,
					toToString,
					valueToEtherThenString,
					gasPaid,
					totalValue)
			}
		}

	}
}
func ListenHeaders(conn *ethclient.Client) {
	headers := make(chan *types.Header)
	subscription, err := conn.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		fmt.Errorf(err.Error())
		color.Red("[OPERATION] Listening to headers behaved unexpectedly.Halting the execution.")
	}

	for {
		select {
		case err := <-subscription.Err():
			log.Fatal(err)
			panic("Something went wrong in the channel")
		case blockHeader := <-headers:
			t := time.Now()
			color.Green("[BLOCK] %d-%02d-%02dT%02d:%02d:%02d-00:00 Block %d found.\n", t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second(), blockHeader.Number.Uint64())
			color.Cyan("[SIZE] Block size:%s\n", blockHeader.Size().String())
			color.Magenta("[BLOCKINFO] Block hash:%s\n", blockHeader.Hash().String())
			fmt.Println("\t(*)parentHash:", blockHeader.ParentHash.Hex())
			fmt.Println("\t(*)unclesHash:", blockHeader.UncleHash.Hex())
			fmt.Println("\t(*)miner:", blockHeader.Coinbase.Hex())
			fmt.Println("\t(*)stateRoot:", blockHeader.Root.Hex())
			fmt.Println("\t(*)transactionsRoot:", blockHeader.TxHash.Hex())
			fmt.Println("\t(*)receiptsRoot:", blockHeader.ReceiptHash.Hex())
			fmt.Println("\t(*)Difficulty:", blockHeader.Difficulty.String())
			fmt.Println("\t(*)GasLimit:", blockHeader.GasLimit)
			fmt.Println("\t(*)Gas Used:", blockHeader.GasUsed)
			fmt.Println("\t(*)Timestamp:", blockHeader.Time)
		}
	}
}
