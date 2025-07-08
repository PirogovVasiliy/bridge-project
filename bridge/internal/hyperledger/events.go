package hyperledger

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type TransferEvent struct {
	Address string `json:"to"`
	Amount  string `json:"value"`
}

func (TrEv TransferEvent) GetAddress() string { return TrEv.Address }
func (TrEv TransferEvent) GetAmount() string  { return TrEv.Amount }

func ListenTransfer(ctx context.Context, network *client.Network, outChanal chan<- TransferEvent) {
	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		log.Fatalln("Ошибка прослушивания событий hyperledger!", err)
	}

	for event := range events {
		var trEvent TransferEvent
		err := json.Unmarshal(event.Payload, &trEvent)
		if err != nil {
			log.Println("Ошибка получения события hyperledger!")
		}
		outChanal <- TransferEvent{trEvent.Address, trEvent.Amount}
	}
}
