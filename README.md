# zkWasm-service-helper for Go

This is the `zkWasm-service-helper` SDK for the Go programming language.

It requires a minimum version of `Go 1.20`.

This library provides `ZkWasmServiceHelper` to help user to communicate to both ZkWasm service backend and smart contract. It mainly provides API to add tasks into ZkWasm or gather information from ZkWasm.

The JS vesion can be found at [DelphinusLab](https://github.com/DelphinusLab/zkWasm-service-helper)

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

This lib main provide a ZkWasmServiceHelper class to help user to communicate to zkwasm service backend. It mainly provide API to add tasks and get informations to zkwasm service backend.

like:

- async queryImage(md5: string);
- async loadTasks(query: QueryParams);
- async addProvingTask(task: WithSignature);
- async ContractVerify(proof []byte, verifyInstance []byte, aux []byte, targetInstance []byte);

### A example to query wasm image
```
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
```

### A example to add new proving task
```
func SubmitProof(ctx context.Context, appID string, publicInputs []any, privateInputs []any) (string, error) {
	appInfo := code.CodeDB[appID]
	if appInfo == nil {
		return "", fmt.Errorf("appID %s not found", appID)
	}

	provingParams := &zkwasm.ProvingParams{
		UserAddress:      zkWasmHelper.GetUserAddress(),
		Md5:              appInfo.WasmMD5,
		InputContextType: zkwasm.ProvingParamsInputContextTypeImageCurrent,
	}

	if len(publicInputs) != 0 {
		pubIArr, err := zkwasm.BuildInputsString(publicInputs)
		if err != nil {
			return "", err
		}
		provingParams.PublicInputs = pubIArr
	}

	if len(privateInputs) != 0 {
		privIArr, err := zkwasm.BuildInputsString(privateInputs)
		if err != nil {
			return "", err
		}
		provingParams.PrivateInputs = privIArr
	}

	id, err := zkWasmHelper.AddProvingTask(ctx, provingParams)
	if err != nil {
		return "", err
	}

	fmt.Printf("task id: %s\n", id)
```

### More examples can be found under `example` project.
