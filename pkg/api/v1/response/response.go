package response

// ErrorResponse represents a standard error response
type ErrorResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details T      `json:"details,omitempty"`
}

// SuccessResponse represents a standard success response
type SuccessResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse[T any](code int, message string, details T) *ErrorResponse[T] {
	return &ErrorResponse[T]{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse[T any](code int, message string, data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Common error messages
const (
	ErrInvalidRequest     = "Invalid request"
	ErrUnauthorized       = "Unauthorized"
	ErrForbidden          = "Forbidden"
	ErrNotFound           = "Not found"
	ErrInternalServer     = "Internal server error"
	ErrValidationFailed   = "Validation failed"
	ErrDuplicateEntry     = "Duplicate entry"
	ErrInvalidToken       = "Invalid token"
	ErrTokenExpired       = "Token expired"
	ErrTokenRevoked       = "Token has been revoked"
	ErrInvalidCredentials = "Invalid credentials"
)
