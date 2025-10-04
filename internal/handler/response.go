package handler

// StatusResponse represents a standard status payload.
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents an error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}
