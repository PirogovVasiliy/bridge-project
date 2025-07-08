package main

import (
	"bridge/internal/bridge"
	"bridge/internal/hardhat"
	"bridge/internal/hyperledger"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	hyperledgerNetwork, contract := hyperledger.ConnectToContract()
	fmt.Println("--- Подключено к hyperledger! ---")
	fmt.Println("--------------------------------------")

	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic(err)
	}

	var networkFST hardhat.Network
	preparation("URL", "DEPLOYER_PRIVATE_KEY", &networkFST)
	fmt.Println("--- Подключено к hardhat! ---")
	fmt.Println("--------------------------------------")

	bridge.Bridge(&networkFST, contract, hyperledgerNetwork)
}

func preparation(node string, pk string, network *hardhat.Network) {
	nodeURL := os.Getenv(node)
	privateKey := os.Getenv(pk)
	err := network.Init(nodeURL, privateKey)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("contract address:", network.GetAddress())
}
