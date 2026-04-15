package types

type StandardResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}
