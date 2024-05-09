package zkwasm

import (
	"context"
	"crypto/ecdsa"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	endpointImage       = "/image"
	endpointImageBinary = "/imagebinary"
	endpointTasks       = "/tasks"
	endpointProve       = "/prove"
	endpointSetup       = "/setup"

	headerSignatureKey = "x-eth-signature"
)

type ZkWasmServiceHelper struct {
	zkWasmEndpoint string
	ethEndpoint    string

	wallet      *ecdsa.PrivateKey
	userAddress common.Address

	ethClient *ethclient.Client

	verifyContractAddress common.Address
}

func New(zkWasmEndpoint, ethEndpoint, privateKey, contractAddress string) (*ZkWasmServiceHelper, error) {
	return NewWithContext(context.Background(), zkWasmEndpoint, ethEndpoint, privateKey, contractAddress)
}

func NewWithContext(ctx context.Context, zkWasmEndpoint, ethEndpoint, privateKey, contractAddress string) (*ZkWasmServiceHelper, error) {
	h := &ZkWasmServiceHelper{}

	h.zkWasmEndpoint = strings.TrimSuffix(zkWasmEndpoint, "/")
	h.ethEndpoint = strings.TrimSuffix(ethEndpoint, "/")

	ethC, err := ethclient.DialContext(ctx, h.ethEndpoint)
	if err != nil {
		return nil, err
	}
	h.ethClient = ethC

	h.verifyContractAddress = common.HexToAddress(contractAddress)

	privK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	h.wallet = privK
	h.userAddress = crypto.PubkeyToAddress(*h.wallet.Public().(*ecdsa.PublicKey))

	return h, nil
}

func (h *ZkWasmServiceHelper) GetUserAddress() string {
	return h.userAddress.Hex()
}

func (h *ZkWasmServiceHelper) signMessage(message string, legacyV bool) (string, error) {
	hash := accounts.TextHash([]byte(message))

	sign, err := crypto.Sign(hash, h.wallet)
	if err != nil {
		return "", err
	}

	// https://github.com/ethereum/go-ethereum/issues/19751
	if legacyV {
		sign[64] += 27
	}

	return hexutil.Encode(sign), nil
}
