package main

import (
	"context"
	"flag"
	"getherscan-cli/cmd/message"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {

}

func main() {
	networkName := flag.String("network", "mainnet", "a string")
	//company := flag.String("provider", "noprovider", "a string")
	//providerNetworkType := flag.String("providerNetworkType", "wss", "a string")
	operation := flag.String("operation", "lens_block", "a string")
	txHash := flag.String("txhash", "0x0", "a string")
	blocknumber := flag.Uint64("blocknumber", 14146232, "an unsigned integer")
	//TODO:Refactor
	flag.Parse()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		message.SetupFailEnvironmentMessage()
		os.Exit(2)
	}
	value, ok := viper.Get("INFURA_KEY").(string)
	if !ok {
		message.SetupFailEnvironmentMessage()
		os.Exit(2)
	}

	message.NetworkDialNodeMessage()

	networkConfig := []string{"wss://", *networkName, ".infura.io/ws/v3/", value}
	conn, err := ethclient.Dial(strings.Join(networkConfig, ""))
	if err != nil {
		message.NetworkFailToDialMessage()
		os.Exit(2)
	}
	// relative connection test, if we have dialed the wrong endpoint,
	// any function of eth.Client should behave unexpectedly.
	networkId, err := conn.NetworkID(context.Background())
	if networkId == nil || err != nil {
		message.NetworkFailToDialMessage()
		os.Exit(2)
	}

	message.NetworkConnectionSuccessfulMessage()

	switch *operation {
	case "lens_transaction":
		if *txHash == "0x0" {
			message.UserFalseParameterMessage()
			os.Exit(2)
		}
		SingleTransaction(conn, *txHash)
	case "lens_block":
		if *blocknumber < 0 {
			message.UserFalseParameterMessage()
			os.Exit(2)
		}
		SingleBlock(conn, *blocknumber)
	}
}
