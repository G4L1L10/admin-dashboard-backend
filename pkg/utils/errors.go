package utils

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewErrorResponse creates a consistent error response
func NewErrorResponse(message string, details string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Details: details,
	}
}
