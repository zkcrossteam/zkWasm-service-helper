package zkwasm

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"slices"
	"strconv"
	"strings"
)

const (
	ImageMetadataKeysProvePaymentSrc = "ProvePaymentSrc"

	ImageMetadataValsProvePaymentSrcDefault    = "Default"
	ImageMetadataValsProvePaymentSrcCreatorPay = "CreatorPay"

	ImageStatusReceived    = "Received"
	ImageStatusInitialized = "Initialized"
	ImageStatusVerified    = "Verified"

	ImageCircuitSizeDefault = 22
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

type AddImageResult struct {
	ID  string `json:"id"`
	MD5 string `json:"md5"`
}

type AddImageParams struct {
	Name           string `json:"name"`
	Image          []byte `json:"image"`
	ImageMD5       string `json:"image_md5"`
	UserAddress    string `json:"user_address"`
	DescriptionUrl string `json:"description_url"`
	AvatorUrl      string `json:"avator_url"`
	CircuitSize    int64  `json:"circuit_size"`

	MetadataKeys []string `json:"metadata_keys"`
	MetadataVals []string `json:"metadata_vals"`

	InitialContext    []byte `json:"initial_context,omitempty"`
	InitialContextMD5 string `json:"initial_context_md5,omitempty"`
}

func (p *AddImageParams) fillValues(userAddress string) {
	// overwrite user address
	p.UserAddress = strings.ToLower(userAddress)

	if p.Image != nil {
		s := md5.Sum(p.Image)
		p.ImageMD5 = strings.ToLower(hex.EncodeToString(s[:]))
	}

	if p.CircuitSize == 0 {
		p.CircuitSize = ImageCircuitSizeDefault
	}

	if !slices.Contains(p.MetadataKeys, ImageMetadataKeysProvePaymentSrc) {
		p.MetadataKeys = append(p.MetadataKeys, ImageMetadataKeysProvePaymentSrc)
		p.MetadataVals = append(p.MetadataVals, ImageMetadataValsProvePaymentSrcDefault)
	}

	if p.InitialContext != nil {
		s := md5.Sum(p.InitialContext)
		p.InitialContextMD5 = strings.ToLower(hex.EncodeToString(s[:]))
	}
}

func (p *AddImageParams) buildSignMessage() string {
	message := ""

	message += p.Name
	message += p.ImageMD5
	message += p.UserAddress
	message += p.DescriptionUrl
	message += p.AvatorUrl
	message += strconv.FormatInt(p.CircuitSize, 10)

	if len(p.MetadataKeys) > 0 {
		message += strings.Join(p.MetadataKeys, ",")
	}
	if len(p.MetadataVals) > 0 {
		message += strings.Join(p.MetadataVals, ",")
	}

	if p.InitialContext != nil {
		message += p.InitialContextMD5
	}

	return message
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

func (h *ZkWasmServiceHelper) AddNewWasmImage(ctx context.Context, params *AddImageParams) (string, error) {
	params.fillValues(h.GetUserAddress())

	signMsg := params.buildSignMessage()
	sign, err := h.signMessage(signMsg, true)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	w.WriteField("name", params.Name)

	imageHeader := make(textproto.MIMEHeader)
	imageHeader.Set("Content-Type", "application/wasm")
	imageHeader.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="image"; filename="%s"`, params.Name))
	iw, err := w.CreatePart(imageHeader)
	if err != nil {
		return "", err
	}
	if _, err := iw.Write(params.Image); err != nil {
		return "", err
	}
	w.WriteField("image_md5", params.ImageMD5)

	w.WriteField("user_address", params.UserAddress)
	w.WriteField("description_url", params.DescriptionUrl)
	w.WriteField("avator_url", params.AvatorUrl)
	w.WriteField("circuit_size", strconv.FormatInt(params.CircuitSize, 10))
	for _, k := range params.MetadataKeys {
		w.WriteField("metadata_keys", k)
	}
	for _, v := range params.MetadataVals {
		w.WriteField("metadata_vals", v)
	}

	if len(params.InitialContext) != 0 {
		icw, err := w.CreateFormField("initial_context")
		if err != nil {
			return "", err
		}

		if _, err := icw.Write(params.InitialContext); err != nil {
			return "", err
		}

		w.WriteField("initial_context_md5", params.InitialContextMD5)
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.zkWasmEndpoint+endpointSetup, &b)
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

	result := &Response[*AddImageResult]{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if !result.Success || result.Result == nil {
		return "", errors.New(string(body))
	}

	return result.Result.ID, nil
}
