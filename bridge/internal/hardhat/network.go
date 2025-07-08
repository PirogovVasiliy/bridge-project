package hardhat

import (
	"bridge/internal/hardhat/contract"
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Network struct {
	nodeURL         string
	ownerPrivateKey *ecdsa.PrivateKey
	client          *ethclient.Client
	chainID         *big.Int
	contract        *contract.BridgeToken
	contractAddress common.Address
}

func (network Network) GetAddress() string { return network.contractAddress.Hex() }

func (network Network) GetContract() *contract.BridgeToken { return network.contract }

func (network Network) GetPK() *ecdsa.PrivateKey { return network.ownerPrivateKey }

func (network Network) GetClient() *ethclient.Client { return network.client }

func (network Network) GetChainID() *big.Int { return network.chainID }

func (network *Network) Init(nodeURL string, privateKeyStr string) (err error) {
	network.nodeURL = nodeURL
	network.ownerPrivateKey, err = crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return err
	}

	network.client, err = connect(network.nodeURL)
	if err != nil {
		return err
	}

	network.chainID, err = network.client.ChainID(context.Background())
	if err != nil {
		return err
	}

	network.contract, network.contractAddress, err = deployContract(network.client, network.ownerPrivateKey, network.chainID)
	if err != nil {
		return err
	}

	return err
}

func connect(nodeURL string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(nodeURL)
	return client, err
}

func deployContract(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
) (btInstance *contract.BridgeToken,
	contractAddress common.Address,
	err error) {

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return
	}

	contractAddress, tx, btInstance, err := contract.DeployBridgeToken(auth, client)
	if err != nil {
		return
	}

	_, err = bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		return
	}
	return
}
