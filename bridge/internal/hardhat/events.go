package hardhat

import (
	"bridge/internal/hardhat/contract"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type TransferEvent struct {
	address string
	amount  *big.Int
}

func (TrEv TransferEvent) GetAddress() string  { return TrEv.address }
func (TrEv TransferEvent) GetAmount() *big.Int { return TrEv.amount }

func ListenTransfer(btInstance *contract.BridgeToken, outChanal chan<- TransferEvent) {

	transferChanal := make(chan *contract.BridgeTokenTokensSentToBridge)

	subscription, err := btInstance.WatchTokensSentToBridge(&bind.WatchOpts{}, transferChanal)
	if err != nil {
		fmt.Println("Ошибка Подписки")
		log.Panic(err)
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err := <-subscription.Err():
			log.Fatalln("Ошибка чтения события TokensSentToBridge!", err)
		case event := <-transferChanal:
			outChanal <- TransferEvent{event.To, event.Amount}
		}
	}
}
