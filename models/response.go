package models

type ResponseWithMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// If you're using Go 1.18+, make it generic:
type ResponseWithData[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}
