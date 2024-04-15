package zkwasm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (h *ZkWasmServiceHelper) QueryImageBinary(ctx context.Context, md5 string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.zkWasmEndpoint+endpointImageBinary, nil)
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

	response := &Response[[]byte]{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, errors.New(err.Error() + ": " + string(body))
	}

	if !response.Success {
		return nil, errors.New(string(body))
	}

	return response.Result, nil
}
