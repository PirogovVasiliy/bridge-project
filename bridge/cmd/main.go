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

const (
	mspID         = "Org1MSP"
	cryptoPath    = "../../network-two/base-network/organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/Admin@org1.example.com/msp/signcerts"
	keyPath       = cryptoPath + "/users/Admin@org1.example.com/msp/keystore"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "dns:///localhost:7151"
	gatewayPeer   = "peer0.org1.example.com"
	chaincodeName = "basic"
	channelName   = "mychannel"
)

func main() {

	hyperledgerNetwork, contract := hyperledger.ConnectToContract(
		mspID, certPath, keyPath, tlsCertPath, peerEndpoint, gatewayPeer, chaincodeName, channelName)
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
