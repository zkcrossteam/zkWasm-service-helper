package main

import (
	"context"
	"fmt"

	helper "github.com/zkcrossteam/zkWasm-service-helper"
)

const (
	zkWasmEndpoint     = "https://rpc.zkwasmhub.com:8090"
	ethEndpoint        = "https://rpc.goerli.eth.gateway.fm"
	privateKey         = "8644db7d9d8beb607960dc23d260d5ac66e8534c41ae77b4d6e22de613d3da2f"
	zkWasmContractAddr = "0x9D48Dce80682108864F1FB719229DCd0C45E51D7"
)

func main() {
	ctx := context.Background()

	h, err := helper.NewWithContext(ctx, zkWasmEndpoint, ethEndpoint, privateKey, zkWasmContractAddr)
	if err != nil {
		panic(err)
	}

	imageMD5 := "FBE1ADD84935782493030FF335475D81"

	i, err := h.QueryImage(ctx, imageMD5)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", i)
}
