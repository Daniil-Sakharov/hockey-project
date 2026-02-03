package dto

// HealthResponse represents health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
