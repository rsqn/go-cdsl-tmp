package exceptions

import (
	"fmt"
)

// CdslError is the base error type for CDSL errors
type CdslError struct {
	Message string
	Cause   error
}

// Error implements the error interface
func (e *CdslError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause
func (e *CdslError) Unwrap() error {
	return e.Cause
}

// NewCdslError creates a new CdslError
func NewCdslError(message string, cause error) *CdslError {
	return &CdslError{
		Message: message,
		Cause:   cause,
	}
}

// CdslValidationError represents a validation error
type CdslValidationError struct {
	CdslError
}

// NewCdslValidationError creates a new CdslValidationError
func NewCdslValidationError(message string, cause error) *CdslValidationError {
	return &CdslValidationError{
		CdslError: CdslError{
			Message: message,
			Cause:   cause,
		},
	}
}
