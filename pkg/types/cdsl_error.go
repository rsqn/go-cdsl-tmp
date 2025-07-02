package types

// CdslError represents an error in the CDSL framework
type CdslError struct {
	Message string
	Cause   error
}

// Error implements the error interface
func (e *CdslError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *CdslError) Unwrap() error {
	return e.Cause
}
