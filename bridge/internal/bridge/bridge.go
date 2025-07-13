package bridge

import (
	"bridge/internal/hardhat"
	"bridge/internal/hyperledger"
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const ownerAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

func Bridge(network *hardhat.Network,
	contractFst *client.Contract, hyperledgerNetworkFst *client.Network,
	contractSec *client.Contract, hyperledgerNetworkSec *client.Network,
) {
	hardChan := make(chan hardhat.TransferEvent)
	go hardhat.ListenTransfer(network.GetContract(), hardChan)
	fmt.Println("Начинаем слушать события с HardHat!")
	fmt.Println("-----------------------------------------")

	eventChanFst := make(chan hyperledger.TransferEvent)
	ctxFst, cancelFst := context.WithCancel(context.Background())
	defer cancelFst()
	eventChanSec := make(chan hyperledger.TransferEvent)
	ctxSec, cancelSec := context.WithCancel(context.Background())
	defer cancelSec()

	go hyperledger.ListenTransfer(ctxFst, hyperledgerNetworkFst, eventChanFst, "basic")
	fmt.Println("Начинаем слушать события первой сети hyperledger!")
	fmt.Println("-----------------------------------------")

	go hyperledger.ListenTransfer(ctxSec, hyperledgerNetworkSec, eventChanSec, "basic")
	fmt.Println("Начинаем слушать события второй сети hyperledger!")
	fmt.Println("-----------------------------------------")

	for {
		select {
		case event := <-hardChan:
			fmt.Println("--- Получено событие с HardHat! ---")
			fmt.Println("Address:", event.GetAddress())
			fmt.Println("Amount:", event.GetAmount())
			fmt.Println("ChainID:", event.GetChainID())

			if event.GetChainID() == 1 {
				transferHyperledgerHelper(contractFst, &event)
			} else if event.GetChainID() == 2 {
				transferHyperledgerHelper(contractSec, &event)
			} else {
				log.Println("Несуществующий ChainID")
			}
		case event := <-eventChanFst:
			fmt.Println("--- Получено событие с первой сети Hyperledger! ---")
			fmt.Println("Address:", event.Address)
			fmt.Println("Amount:", event.Amount)
			fmt.Println("-----------------------------------------")
			transferHardhatHelper(network, &event, 2)
		case event := <-eventChanSec:
			fmt.Println("--- Получено событие со второй сети Hyperledger! ---")
			fmt.Println("Address:", event.Address)
			fmt.Println("Amount:", event.Amount)
			fmt.Println("-----------------------------------------")
			transferHardhatHelper(network, &event, 1)
		}
	}
}

func transferHyperledgerHelper(contract *client.Contract, event *hardhat.TransferEvent) {
	err := hyperledger.CallMint(contract, event.GetAmount().String())
	if err != nil {
		log.Fatalln("Ошибка вызова Mint")
	}
	err = hyperledger.CallTransfer(contract, event.GetAddress(), event.GetAmount().String())
	if err != nil {
		log.Fatalln("Ошибка вызова Transfer")
	}

	fmt.Println("Перевод прошел успешно!")
	fmt.Println("-----------------------------------------")
}

func transferHardhatHelper(hardhatNetwork *hardhat.Network, event *hyperledger.TransferEvent, chainID int) {
	err := hardhat.CallReceiveFromBridge(
		hardhatNetwork.GetClient(),
		hardhatNetwork.GetContract(),
		hardhatNetwork.GetPK(),
		hardhatNetwork.GetChainID(),
		common.HexToAddress(ownerAddress),
		event.Amount,
	)
	if err != nil {
		log.Fatalln("Ошибка вызова ReceiveFromBridge!", err)
	}

	err = hardhat.CallSendToBridge(
		hardhatNetwork.GetClient(),
		hardhatNetwork.GetContract(),
		hardhatNetwork.GetPK(),
		hardhatNetwork.GetChainID(),
		event.GetAddress(),
		event.GetAmount(),
		chainID,
	)
	if err != nil {
		log.Fatalln("Ошибка вызова SendToBridge!", err)
	}
}
