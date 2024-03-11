package zkwasm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	verifyABIJSON = `
[
  {
    "inputs":
    [
      {
        "internalType": "uint256[]",
        "name": "proof",
        "type": "uint256[]"
      },
      {
        "internalType": "uint256[]",
        "name": "verify_instance",
        "type": "uint256[]"
      },
      {
        "internalType": "uint256[]",
        "name": "aux",
        "type": "uint256[]"
      },
      {
        "internalType": "uint256[][]",
        "name": "target_instance",
        "type": "uint256[][]"
      }
    ],
    "name": "verify",
    "outputs": [],
    "stateMutability": "view",
    "type": "function"
  }
]
`
)

func (h *ZkWasmServiceHelper) ContractVerify(ctx context.Context,
	proof []byte, verifyInstance []byte, aux []byte, targetInstance []byte) (string, error) {

	verifyABI, err := abi.JSON(strings.NewReader(verifyABIJSON))
	if err != nil {
		return "", err
	}

	data, err := verifyABI.Pack("verify", ByteSliceToBigIntSlice(proof, true), ByteSliceToBigIntSlice(verifyInstance, true), ByteSliceToBigIntSlice(aux, true), [][]*big.Int{ByteSliceToBigIntSlice(targetInstance, true)})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	chainID, err := h.ethClient.ChainID(ctx)
	if err != nil {
		return "", err
	}

	nonce, err := h.ethClient.PendingNonceAt(ctx, h.userAddress)
	if err != nil {
		return "", err
	}

	gasPrice, err := h.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	gasTipCap, err := h.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return "", err
	}

	signTx, err := types.SignNewTx(h.wallet, types.NewLondonSigner(chainID), &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasPrice,
		Gas:       uint64(1500000),
		To:        &h.verifyContractAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	err = h.ethClient.SendTransaction(ctx, signTx)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return signTx.Hash().Hex(), nil
}
