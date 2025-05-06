package customErrors

type InternalError struct {
	Message string `json:"message"`
}

func NewInternal(message string) *InternalError {
	return &InternalError{
		Message: message,
	}
}

func (e *InternalError) Error() string {
	return "InternalError: " + e.Message
}
