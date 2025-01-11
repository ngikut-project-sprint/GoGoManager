package utils

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

type ErrorType int

const (
	InvalidEmailFormat ErrorType = iota
	InvalidURIFormat
	InvalidPasswordLength
	InvalidUserId
	InvalidNameLength
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

func NoError(t assert.TestingT, err *GoGoError, msgAndArgs ...interface{}) bool {
	var e error
	if err != nil {
		e = err.Err
	} else {
		e = nil
	}
	return assert.NoError(t, e)
}

func Error(t assert.TestingT, err *GoGoError, msgAndArgs ...interface{}) bool {
	var e error
	if err != nil {
		e = err.Err
	} else {
		e = nil
	}
	return assert.Error(t, e)
}
