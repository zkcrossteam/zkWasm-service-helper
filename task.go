package zkwasm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

const (
	TaskStatusPending       = "Pending"
	TaskStatusProcessing    = "Processing"
	TaskStatusDryRunSuccess = "DryRunSuccess"
	TaskStatusDryRunFailed  = "DryRunFailed"
	TaskStatusDone          = "Done"
	TaskStatusFail          = "Fail"
	TaskStatusStale         = "Stale"
)

type TaskQueryParams struct {
	UserAddress string `json:",omitempty"`
	MD5         string `json:"md5,omitempty"`
	ID          string `json:"id,omitempty"`
	TaskType    string `json:"tasktype,omitempty"`
	TaskStatus  string `json:"taskstatus,omitempty"`
	Start       int64  `json:"start,omitempty"`
	Total       int64  `json:"total,omitempty"`
}

type Task struct {
	UserAddress       string   `json:"user_address"`
	NodeAddress       string   `json:"node_address"`
	MD5               string   `json:"md5,omitempty"`
	TaskType          string   `json:"task_type,omitempty"`
	Status            string   `json:"status"`
	SingleProof       []byte   `json:"single_proof"`
	Proof             []byte   `json:"proof"`
	Aux               []byte   `json:"aux"`
	ExternalHostTable []byte   `json:"external_host_table"`
	ShadowInstances   []byte   `json:"shadow_instances"`
	BatchInstances    []byte   `json:"batch_instances"`
	Instances         []byte   `json:"instances"`
	PublicInputs      []string `json:"public_inputs"`
	PrivateInputs     []string `json:"private_inputs"`
	InputContext      []byte
	InputContextType  string
	OutputContext     []byte
	ID                string `json:"id"`
	SubmitTime        string
	ProcessStarted    string
	ProcessFinished   string
	TaskFee           []byte `json:"task_fee"`
	StatusMessage     string
	InternalMessage   string
	// task_verification_data: TaskVerificationData;
	DebugLogs       string
	ProofSubmitMode string `json:"proof_submit_mode"`
	// batch_proof_data?: BatchProofData;
	AutoSubmitStatus string `json:"auto_submit_status"`
}

func (h *ZkWasmServiceHelper) LoadTasks(ctx context.Context, query *TaskQueryParams) (*PaginationResult[*Task], error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.zkWasmEndpoint+endpointTasks, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if query.ID != "" {
		q.Add("id", query.ID)
	}
	if query.Total != 0 {
		q.Add("total", strconv.FormatInt(query.Total, 10))
	}
	// TODO: other params

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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	result := &Response[*PaginationResult[*Task]]{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	if !result.Success || result.Result == nil {
		return nil, errors.New(string(body))
	}

	return result.Result, nil
}
