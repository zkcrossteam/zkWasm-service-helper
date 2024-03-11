package zkwasm

type Response[T any] struct {
	Success bool `json:"success"`
	Result  T    `json:"result"`
}

type PaginationResult[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}
