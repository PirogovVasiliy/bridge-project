package hardhat

import (
	"bridge/internal/hardhat/contract"
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CallBuyTokens(
	client *ethclient.Client,
	btInstance *contract.BridgeToken,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
	amount *big.Int,
) error {

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}
	auth.Value = amount

	tx, err := btInstance.BuyTokens(auth)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), client, tx)

	return err
}

func CallReceiveFromBridge(
	client *ethclient.Client,
	btInstance *contract.BridgeToken,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
	to common.Address,
	amountStr string,
) error {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}

	amount := big.NewInt(0)
	amount.SetString(amountStr, 10)
	tx, err := btInstance.ReceiveFromBridge(auth, to, amount)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return err
	}
	return nil
}
