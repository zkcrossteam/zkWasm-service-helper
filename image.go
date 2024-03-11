package zkwasm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Image struct {
	UserAddress string `json:"user_address"`
	MD5         string `json:"md5"`
	// deployment: Array<DeploymentInfo>;
	DescriptionUrl string `json:"description_url"`
	AvatorUrl      string `json:"avator_url"`
	CircuitSize    int64  `json:"circuit_size"`
	Context        []byte `json:"context"`
	InitialContext []byte `json:"initial_context"`
	Status         string `json:"status"`
	// checksum: ImageChecksum | null;
}

func (h *ZkWasmServiceHelper) QueryImage(ctx context.Context, md5 string) (*Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.zkWasmEndpoint+endpointImage, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("md5", md5)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response[[]*Image]{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(string(body))
	}

	if len(response.Result) == 0 {
		return nil, nil
	}

	return response.Result[0], nil
}
