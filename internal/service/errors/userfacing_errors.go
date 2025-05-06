package customErrors

type UserFacingError struct {
	Message string `json:"message"`
}

func NewUserFacing(message string) *UserFacingError {
	return &UserFacingError{
		Message: message,
	}
}

func (e *UserFacingError) Error() string {
	return e.Message
}
