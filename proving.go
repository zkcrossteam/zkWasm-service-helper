package zkwasm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	ProvingParamsInputContextTypeImageCurrent = "ImageCurrent"
	ProvingParamsInputContextTypeCustom       = "Custom"
)

type ProvingParams struct {
	UserAddress string
	MD5         string

	PublicInputs  []string `json:",omitempty"`
	PrivateInputs []string `json:",omitempty"`

	InputContextType string    `json:"input_context_type,omitempty"`
	InputContext     io.Reader `json:"input_context,omitempty"`
	InputContextMD5  string    `json:"input_context_md5,omitempty"`
}

func (p *ProvingParams) buildSignMessage() string {
	message := ""

	message += p.UserAddress
	message += p.MD5

	for _, i := range p.PublicInputs {
		message += i
	}

	for _, i := range p.PrivateInputs {
		message += i
	}

	if p.InputContextType == ProvingParamsInputContextTypeCustom && p.InputContext != nil {
		message += p.InputContextMD5
	}

	if p.InputContextType != "" {
		message += p.InputContextType
	}

	return message
}

type ProvingResult struct {
	MD5 string `json:"md5"`
	ID  string `json:"id"`
}

func (h *ZkWasmServiceHelper) AddProvingTask(ctx context.Context, params *ProvingParams) (string, error) {
	signMsg := params.buildSignMessage()
	sign, err := h.signMessage(signMsg, false)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("user_address", params.UserAddress)
	w.WriteField("md5", params.MD5)
	for _, i := range params.PublicInputs {
		w.WriteField("public_inputs", i)
	}
	for _, i := range params.PrivateInputs {
		w.WriteField("private_inputs", i)
	}
	if params.InputContextType != "" {
		w.WriteField("input_context_type", params.InputContextType)
	}
	if params.InputContextMD5 != "" {
		w.WriteField("input_context_md5", params.InputContextMD5)
	}
	if params.InputContext != nil {
		nw, err := w.CreateFormField("input_context")
		if err != nil {
			return "", err
		}
		_, err = io.Copy(nw, params.InputContext)
		if err != nil {
			return "", err
		}
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.zkWasmEndpoint+endpointProve, &b)
	if err != nil {
		return "", err
	}

	req.Header[headerSignatureKey] = []string{sign}
	req.Header.Add("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	result := &Response[*ProvingResult]{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if !result.Success || result.Result == nil {
		return "", errors.New(string(body))
	}

	return result.Result.ID, nil
}
