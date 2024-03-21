# zkWasm-service-helper for Go

This is the `zkWasm-service-helper` SDK for the Go programming language.

It requires a minimum version of `Go 1.20`.

This library provides `ZkWasmServiceHelper` to help user to communicate to both ZkWasm service backend and smart contract. It mainly provides API to add tasks into ZkWasm or gather information from ZkWasm.

## Getting started

###### Add Dependencies
```sh
$ go get -u github.com/zkcrossteam/zkWasm-service-helper
```

###### Write Code
In your preferred editor add the following content to `main.go`

```go
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

	fmt.Printf("connected. user address: %s\n", h.GetUserAddress())
}
```

###### Compile and Execute
```sh
$ go run .
connected. user address: 0x582867abf63EbfA5803327a960AC381C2E01d67b
```

## Examples
more examples can be found under `example` project.
