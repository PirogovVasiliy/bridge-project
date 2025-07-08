package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

const nameKey = "name"
const symbolKey = "symbol"

type eventSentToBridge struct {
	To    string `json:"to"`
	Value string `json:"value"`
}

func (s *SmartContract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string) (bool, error) {

	bytes, err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	err = ctx.GetStub().PutState(nameKey, []byte(name))
	if err != nil {
		return false, fmt.Errorf("failed to set token name: %v", err)
	}

	err = ctx.GetStub().PutState(symbolKey, []byte(symbol))
	if err != nil {
		return false, fmt.Errorf("failed to set symbol: %v", err)
	}

	return true, nil
}

func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, amountStr string) error {

	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return errors.New("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	amount := big.NewInt(0)
	amount.SetString(amountStr, 10)

	if amount.Cmp(big.NewInt(0)) < 0 {
		return errors.New("mint amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance = big.NewInt(0)

	if currentBalanceBytes != nil {
		currentBalance.SetString(string(currentBalanceBytes), 10)
	}

	updatedBalance := big.NewInt(0)
	updatedBalance.Add(currentBalance, amount)

	err = ctx.GetStub().PutState(minter, []byte(updatedBalance.String()))
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, amountStr string) error {

	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return errors.New("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	amount := big.NewInt(0)
	amount.SetString(amountStr, 10)

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("mint amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance = big.NewInt(0)

	if currentBalanceBytes != nil {
		currentBalance.SetString(string(currentBalanceBytes), 10)
	}

	updatedBalance := big.NewInt(0)
	updatedBalance.Sub(currentBalance, amount)

	err = ctx.GetStub().PutState(minter, []byte(updatedBalance.String()))
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, recipient string, amountStr string) error {

	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	err = transferHelper(ctx, clientID, recipient, amountStr)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}

	return nil
}

func (s *SmartContract) SendTokensToBridge(ctx contractapi.TransactionContextInterface, recipient string, amountStr string) error {
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	err = s.Burn(ctx, amountStr)
	if err != nil {
		return err
	}

	event := eventSentToBridge{recipient, amountStr}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = ctx.GetStub().SetEvent("SendTokensToBridge", eventJSON)
	return err
}

func (s *SmartContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (string, error) {

	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "0", fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return "0", fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "0", fmt.Errorf("failed to get client id: %v", err)
	}

	balanceBytes, err := ctx.GetStub().GetState(clientID)
	if err != nil {
		return "0", fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return "0", fmt.Errorf("the account %s does not exist", clientID)
	}

	balance := string(balanceBytes)

	return balance, nil
}

func (s *SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {
	initialized, err := checkInitialized(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to check if contract is already initialized: %v", err)
	}
	if !initialized {
		return "", fmt.Errorf("contract options need to be set before calling any function, call Initialize() to initialize contract")
	}

	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}

	return clientAccountID, nil
}

func transferHelper(ctx contractapi.TransactionContextInterface, from string, to string, valueStr string) error {

	if from == to {
		return fmt.Errorf("cannot transfer to and from same client account")
	}

	value := big.NewInt(0)
	value.SetString(valueStr, 10)

	if value.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("transfer amount cannot be negative")
	}

	fromCurrentBalanceBytes, err := ctx.GetStub().GetState(from)
	if err != nil {
		return fmt.Errorf("failed to read client account %s from world state: %v", from, err)
	}

	if fromCurrentBalanceBytes == nil {
		return fmt.Errorf("client account %s has no balance", from)
	}

	var fromCurrentBalance = big.NewInt(0)
	fromCurrentBalance.SetString(string(fromCurrentBalanceBytes), 10)

	if fromCurrentBalance.Cmp(value) < 0 {
		return fmt.Errorf("client account %s has insufficient funds", from)
	}

	toCurrentBalanceBytes, err := ctx.GetStub().GetState(to)
	if err != nil {
		return fmt.Errorf("failed to read recipient account %s from world state: %v", to, err)
	}

	toCurrentBalance := big.NewInt(0)

	if toCurrentBalanceBytes != nil {
		toCurrentBalance.SetString(string(toCurrentBalanceBytes), 10)
	}

	fromUpdatedBalance := big.NewInt(0)
	fromUpdatedBalance.Sub(fromCurrentBalance, value)

	toUpdatedBalance := big.NewInt(0)
	toUpdatedBalance.Add(toCurrentBalance, value)

	err = ctx.GetStub().PutState(from, []byte(fromUpdatedBalance.String()))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(to, []byte(toUpdatedBalance.String()))
	if err != nil {
		return err
	}

	return nil
}

func checkInitialized(ctx contractapi.TransactionContextInterface) (bool, error) {
	tokenName, err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("failed to get token name: %v", err)
	}

	if tokenName == nil {
		return false, nil
	}

	return true, nil
}
