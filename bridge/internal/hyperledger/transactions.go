package hyperledger

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func CallInitialize(contract *client.Contract, name string, symbol string) error {
	_, err := contract.SubmitTransaction("Initialize", name, symbol)
	if err != nil {
		return err
	}
	return nil
}

func CallMint(contract *client.Contract, amount string) error {
	_, err := contract.SubmitTransaction("Mint", amount)
	if err != nil {
		return err
	}
	return nil
}

func CallClientAccountBalance(contract *client.Contract) error {
	_, err := contract.EvaluateTransaction("ClientAccountBalance")
	if err != nil {
		return err
	}
	return nil
}

func CallBurn(contract *client.Contract, amount string) error {
	_, err := contract.SubmitTransaction("Burn", amount)
	if err != nil {
		return err
	}
	return nil
}

func CallClientAccountID(contract *client.Contract) error {
	_, err := contract.EvaluateTransaction("ClientAccountID")
	if err != nil {
		return err
	}
	return nil
}

func CallTransfer(contract *client.Contract, recipient string, amount string) error {
	_, err := contract.SubmitTransaction("Transfer", recipient, amount)
	if err != nil {
		return err
	}
	return nil
}
