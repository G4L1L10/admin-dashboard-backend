package utils

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewSuccessResponse creates a consistent success response
func NewSuccessResponse(message string, data any) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}
