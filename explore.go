package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/status-im/keycard-go/hexutils"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type PendingTxMessage struct {
	Version string
	Method  string
	Params  ParamsResult
}

type ParamsResult struct {
	Subscription string
	Result       string
}

func main() {
	networkNamePtr := flag.String("network", "mainnet", "a string")
	operationPtr := flag.String("operation", "lens_txpool", "a string")
	txHashPtr := flag.String("txhash", "0x0", "a string")
	blockNoPtr := flag.Uint64("blocknumber", 0, "an int")
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
			color.Red("[PARAMETER_ERROR] Tx Hash invalid or not given for this operation.Pass it by --txhash=<transaction hash>.")
		}
		LensTransaction(conn, *txHashPtr)
	case "lens_block":
		if *blockNoPtr < 0 {
			panic("Argument error. There is no block with block number less than 0.")
		}
		LensBlock(conn, *blockNoPtr)
	case "lens_txpool":
		LensTxpool(conn, networkConfig)
	}
}

func LensTxpool(ethconn *ethclient.Client, networkConfig []string) {
	//@TODO:Dirty research. Dig in to the details
	defer func() {
		if err := recover(); err != nil {
			color.Red("[OPERATION] Critical error occurred. Recovering, see the logs.")
			log.Println("Panic:", err)
		}
	}()
	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(strings.Join(networkConfig, ""), nil)
	if err != nil {
		color.Red("[NETWORK] Could not connect to the node directly. Halting the execution.")
		os.Exit(42)
	}
	defer conn.Close()
	done := make(chan struct{})
	var pendingTx PendingTxMessage
	counter := 0 //purpose of this counter is that the first message coming from the channel is unexpected.
	//plus we can use it as the counter for the number of txs in the txpool.
	go func() {
		defer close(done)
		for {
			// Problem: Message is unexpected, solve it.
			_, data, err := conn.ReadMessage()
			if err != nil {
				color.Red("[NETWORK] Could not read the message, or session terminated. If non intended, maybe provider error ? ")
				return
			}
			if counter != 0 {
				json.Unmarshal(data, &pendingTx)
				//@TODO:Fix this
				LensTxOnlyValue(ethconn, pendingTx.Params.Result)

			}
			counter++

		}
	}()
	//new pending tx's
	txPoolRequest := `{"jsonrpc":"2.0", "id": 1, "method": "eth_subscribe", "params": ["newPendingTransactions"]}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(txPoolRequest)); err != nil {
		color.Red("[NETWORK-OPERATION] Could not send  the message. Halting the execution.")
		os.Exit(21)
	}
	for {
		select {
		case <-done:
			return
		case _ = <-messageOut:
			err := conn.WriteMessage(websocket.TextMessage, []byte(txPoolRequest))
			if err != nil {
				color.Red("[NETWORK] Could not write the message.")
				return
			}
		case <-interrupt:
			color.Red("[NETWORK] Interrupt on connection.")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				color.Red("[NETWORK] Could not close the channel")
				return
			}
			select {
			case <-done:
			}
			return
		}
	}
}

// Refactor
func LensBlock(conn *ethclient.Client, blockNumber uint64) {
	if blockNumber < 0 {
		color.Red("[USER] There is no such block.")
		panic("blockNumber is negative.")
	}
	block, err := conn.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))

	if err != nil {
		color.Red("[OPERATION] Block could not be fetched. Maybe check block number ? ")
	}
	color.Green("[BLOCK] Block %d found.\n", block.Header().Number.Uint64())
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
func LensTransaction(conn *ethclient.Client, transaction string) {
	txHash := common.HexToHash(transaction)
	tx, isPending, err := conn.TransactionByHash(context.Background(), txHash)
	//@TODO:Can't figure out why even though websocket returns the hash.
	//Simple dirty-fix at the moment.
	if tx == nil {

		color.Red("[OPERATION] Error occurred while fetching the transaction. Maybe check the transaction hash ?")
		return
	}
	if err != nil {
		color.Red("[OPERATION] Error occurred while fetching the transaction. Maybe check the transaction hash ?")
	}

	networkId, err := conn.NetworkID(context.Background())
	if err != nil {
		color.Red("[OPERATION] Network Id could not be fetched.Halting the execution")
	}

	txAsMessage, err := tx.AsMessage(types.NewLondonSigner(networkId), nil)
	if err != nil {
		color.Red("[OPERATION] Formatting as message failed. Halting the execution.")
	}
	
	color.Cyan("[TRANSACTION] Tx Hash: %s Is Pending ?:%t", tx.Hash().String(), isPending)
	color.Green("[FROM]:%s --> [TO]:%s", txAsMessage.From().String(), txAsMessage.To().String())
	color.Blue("Gas Limit:%d | Value:%d ", txAsMessage.Gas(), txAsMessage.Value().Uint64())
	color.Yellow("[DATA] Transaction data:")
	fmt.Println("\t[HEX] =", "0x"+hexutils.BytesToHex(txAsMessage.Data()))

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
		_ = fmt.Errorf(err.Error())
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

func LensTxOnlyValue(conn *ethclient.Client, transaction string) {
	//@TODO: Execute data while waiting, see the expected outcome ?
	//Exclude pure transactions ?
	defer func() {
		if err := recover(); err != nil {
			color.Red("[OPERATION] Critical error occurred. Recovering, see the logs.")
			log.Println("Panic:", err)
		}
	}()
	txHash := common.HexToHash(transaction)
	tx, _, err := conn.TransactionByHash(context.Background(), txHash)
	//@TODO:Can't figure out why even though websocket returns the hash.
	//Simple dirty-fix at the moment.
	if tx == nil {
		// put them in a queue, try to execute until every tx is processed.
		//color.Red("[OPERATION] Error occurred while fetching the transaction. Maybe check the transaction hash ?")
		return
	}
	if err != nil {
		color.Red("[OPERATION] Error occurred while fetching the transaction. Maybe check the transaction hash ?")
	}

	networkId, err := conn.NetworkID(context.Background())
	if err != nil {
		color.Red("[OPERATION] Network Id could not be fetched.Halting the execution")
	}

	txAsMessage, err := tx.AsMessage(types.NewLondonSigner(networkId), nil)
	if err != nil {
		color.Red("[OPERATION] Formatting as message failed. Halting the execution.")
	}
	color.Green("[TX] TxHash:%s", tx.Hash().String())
	color.Yellow("\t[TX] From:%s --> To:%s | Value: %d | Gas: %d ", txAsMessage.From().String(), txAsMessage.To().String(), txAsMessage.Value().Uint64(), txAsMessage.Gas())

	color.Magenta("\t[TX] Data:%s", "0x"+hexutils.BytesToHex(txAsMessage.Data()))

}
