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

func Bridge(network *hardhat.Network,
	contractFst *client.Contract, hyperledgerNetworkFst *client.Network,
	contractSec *client.Contract, hyperledgerNetworkSec *client.Network,
) {
	hardChan := make(chan hardhat.TransferEvent)
	go hardhat.ListenTransfer(network.GetContract(), hardChan)
	fmt.Println("Начинаем слушать события с HardHat!")
	fmt.Println("-----------------------------------------")

	eventChanFst := make(chan hyperledger.TransferEvent)
	eventChanSec := make(chan hyperledger.TransferEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hyperledger.ListenTransfer(ctx, hyperledgerNetworkFst, eventChanFst, "basic")
	fmt.Println("Начинаем слушать события первой сети hyperledger!")
	fmt.Println("-----------------------------------------")

	go hyperledger.ListenTransfer(ctx, hyperledgerNetworkSec, eventChanSec, "basic")
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
			transferHardhatHelper(network, &event)
		case event := <-eventChanSec:
			fmt.Println("--- Получено событие со второй сети Hyperledger! ---")
			fmt.Println("Address:", event.Address)
			fmt.Println("Amount:", event.Amount)
			fmt.Println("-----------------------------------------")
			transferHardhatHelper(network, &event)
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

func transferHardhatHelper(hardhatNetwork *hardhat.Network, event *hyperledger.TransferEvent) {
	err := hardhat.CallReceiveFromBridge(
		hardhatNetwork.GetClient(),
		hardhatNetwork.GetContract(),
		hardhatNetwork.GetPK(),
		hardhatNetwork.GetChainID(),
		common.HexToAddress(event.Address),
		event.Amount,
	)
	if err != nil {
		log.Fatalln("Ошибка вызова ReceiveFromBridge!", err)
	}
}
