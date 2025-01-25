package models

type ErrorResponse struct {
	Status       int64  `json:"status"`
	ErrorMessage string `json:"error_message"`
}
