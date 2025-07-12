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

func Bridge(network *hardhat.Network, contract *client.Contract, hyperledgerNetwork *client.Network) {

	hardChan := make(chan hardhat.TransferEvent)
	go hardhat.ListenTransfer(network.GetContract(), hardChan)
	fmt.Println("Начинаем слушать события с HardHat!")
	fmt.Println("-----------------------------------------")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	hyperChan := make(chan hyperledger.TransferEvent)
	go hyperledger.ListenTransfer(ctx, hyperledgerNetwork, hyperChan, "basic")
	fmt.Println("Начинаем слушать события с Hyperledger!")
	fmt.Println("-----------------------------------------")

	for {
		select {
		case event := <-hardChan:
			fmt.Println("--- Получено событие с HardHat! ---")
			fmt.Println("Address:", event.GetAddress())
			fmt.Println("Amount:", event.GetAmount())

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
		case event := <-hyperChan:
			fmt.Println("--- Получено событие с Hyperledger! ---")
			fmt.Println("Address:", event.Address)
			fmt.Println("Amount:", event.Amount)
			fmt.Println("-----------------------------------------")

			err := hardhat.CallReceiveFromBridge(
				network.GetClient(),
				network.GetContract(),
				network.GetPK(),
				network.GetChainID(),
				common.HexToAddress(event.Address),
				event.Amount,
			)
			if err != nil {
				log.Fatalln("Ошибка вызова ReceiveFromBridge!", err)
			}

			fmt.Println("Перевод прошел успешно!")
			fmt.Println("-----------------------------------------")
		}
	}
}
