package utils

import "fmt"

type ErrorType int

const (
	InvalidEmailFormat ErrorType = iota
	InvalidURIFormat
	InvalidPasswordLength
	InvalidUserId
	SQLError
	SQLUniqueViolated
	PasswordHashFailed
)

type GoGoError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *GoGoError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Message, e.Err)
}

func WrapError(err error, errorType ErrorType, message string) *GoGoError {
	return &GoGoError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}
