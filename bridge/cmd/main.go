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

const cryptoPathFst = "../../hyperledger-fabric/generate-network/organizations/peerOrganizations/logitech.com"
const cryptoPathSec = "../../network-two/base-network/organizations/peerOrganizations/org1.example.com"

type hyperledgerNetwork struct {
	mspID         string
	certPath      string
	keyPath       string
	tlsCertPath   string
	peerEndpoint  string
	gatewayPeer   string
	chaincodeName string
	channelName   string
}

func main() {

	hNetFst := hyperledgerNetwork{
		"LogitechMSP",
		cryptoPathFst + "/users/Admin@logitech.com/msp/signcerts",
		cryptoPathFst + "/users/Admin@logitech.com/msp/keystore",
		cryptoPathFst + "/peers/peer0.logitech.com/tls/ca.crt",
		"dns:///localhost:11051",
		"peer0.logitech.com",
		"basic",
		"mychannel",
	}

	hNetSec := hyperledgerNetwork{
		"Org1MSP",
		cryptoPathSec + "/users/Admin@org1.example.com/msp/signcerts",
		cryptoPathSec + "/users/Admin@org1.example.com/msp/keystore",
		cryptoPathSec + "/peers/peer0.org1.example.com/tls/ca.crt",
		"dns:///localhost:7151",
		"peer0.org1.example.com",
		"basic",
		"mychannel",
	}

	hyperledgerNetworkFst, hyperledgerContractFst := hyperledger.ConnectToContract(
		hNetFst.mspID, hNetFst.certPath, hNetFst.keyPath, hNetFst.tlsCertPath, hNetFst.peerEndpoint,
		hNetFst.gatewayPeer, hNetFst.chaincodeName, hNetFst.channelName)
	fmt.Println("--- Подключено к hyperledger first network! ---")
	fmt.Println("--------------------------------------")

	hyperledgerNetworkSec, hyperledgerContractSec := hyperledger.ConnectToContract(
		hNetSec.mspID, hNetSec.certPath, hNetSec.keyPath, hNetSec.tlsCertPath, hNetSec.peerEndpoint,
		hNetSec.gatewayPeer, hNetSec.chaincodeName, hNetSec.channelName)
	fmt.Println("--- Подключено к hyperledger second network! ---")
	fmt.Println("--------------------------------------")

	err := godotenv.Load("../.env")

	if err != nil {
		log.Panic(err)
	}

	var networkFST hardhat.Network
	preparation("URL", "DEPLOYER_PRIVATE_KEY", &networkFST)
	fmt.Println("--- Подключено к hardhat! ---")
	fmt.Println("--------------------------------------")

	bridge.Bridge(&networkFST,
		hyperledgerContractFst, hyperledgerNetworkFst,
		hyperledgerContractSec, hyperledgerNetworkSec,
	)

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
