package model

type Response[T any] struct {
	Status string `json:"status"`
	Data   *T     `json:"data,omitempty"`
	Errors []any  `json:"errors,omitempty"`
}
