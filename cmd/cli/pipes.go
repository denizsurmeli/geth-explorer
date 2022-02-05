package main

import (
	"getherscan-cli/cmd/message"
	"getherscan-cli/pkg/blockops"
	"getherscan-cli/pkg/transactionops"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
)

// TODO:Refactor the commented lines.
//func ReadProviderConfiguration(filePath string) (*provider.Config, error) {
//	f, err := os.Open(filePath)
//	defer f.Close()
//	if err != nil {
//		message.OperationErrorMessage()
//		return nil, errors.New("Provider configuration read failed.")
//	}
//	byteValue, err := ioutil.ReadAll(f)
//	if err != nil {
//		message.OperationErrorMessage()
//		return nil, errors.New("Provider configuration byte parsing failed")
//	}
//	var config *provider.Config
//
//	err = json.Unmarshal(byteValue, &config)
//	if err != nil {
//		message.OperationErrorMessage()
//		return nil, errors.New("Unmarshal failed.")
//	}
//
//	return config, nil
//}
//
//func Environment(filePath string, company string) (string, error) {
//	viper.SetConfigFile(filePath)
//	err := viper.ReadInConfig()
//	if err != nil {
//		message.SetupFailEnvironmentMessage()
//	}
//	switch company {
//	case "infura":
//		key, ok := viper.Get("INFURA_KEY").(string)
//		if !ok {
//			message.SetupFailEnvironmentMessage()
//		}
//		return key, nil
//	case "alchemy":
//		key, ok := viper.Get("ALCHEMY_KEY").(string)
//		if !ok {
//			message.SetupFailEnvironmentMessage()
//		}
//		return key, nil
//	case "quicknode":
//		key, ok := viper.Get("QUICKNODE_KEY").(string)
//		if !ok {
//			message.SetupFailEnvironmentMessage()
//		}
//		return key, nil
//	default:
//		message.SetupFailEnvironmentMessage()
//		return "", errors.New("Unknown provider.")
//	}
//}
//
//func WrapProvider(configPath string, envPath string, who int, conntype int, net int) (*provider.Provider, error) {
//	conf, err := ReadProviderConfiguration(configPath)
//	if err != nil {
//		message.SetupFailEnvironmentMessage()
//		return nil, errors.New("Could not read config.")
//	}
//	key, err := Environment(envPath)
//	if err != nil {
//		message.SetupFailEnvironmentMessage()
//		return nil, errors.New("Could not read the key.")
//	}
//
//	var p *provider.Provider
//
//}
//
//func Connection(provider provider.Provider) (*ethclient.Client, error) {
//	switch provider.Net {
//	case 2:
//		conn, err := ethclient.Dial(provider.FullUrl)
//		if err != nil {
//			message.NetworkFailToDialMessage()
//			return nil, errors.New("No connection is established.")
//		}
//
//		return conn, nil
//	default:
//		message.NonSupportedOperation()
//		return nil, errors.New("Unsupported operation.")
//	}
//
//}

func SingleTransaction(conn *ethclient.Client, txHex string) {
	txHash := common.HexToHash(txHex)
	err := transactionops.GuardedSingleTransaction(conn, txHash)
	if err != nil {
		message.OperationFailTransaction()
		os.Exit(2)
	}
}

func SingleBlock(conn *ethclient.Client, blocknumber uint64) {
	err := blockops.GuardedSingleBlock(conn, blocknumber)
	if err != nil {
		message.OperationErrorMessage()
		os.Exit(2)
	}
}
