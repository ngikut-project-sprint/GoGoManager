package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

func TestValidateEmail_Valid(t *testing.T) {
	email := "example@domain.com"
	err := validators.ValidateEmail(email)
	assert.NoError(t, err)
}

func TestValidateEmail_WrongFormat(t *testing.T) {
	email := "example@domain.com.abc.invalid@email"
	err := validators.ValidateEmail(email)
	assert.Error(t, err)
}

func TestValidateEmail_WrongRFCFormat(t *testing.T) {
	email := "domain.com"
	err := validators.ValidateEmail(email)
	assert.Error(t, err)
}

func TestValidateEmail_DomainNotFound(t *testing.T) {
	email := "example@domain.c"
	err := validators.ValidateEmail(email)
	assert.Error(t, err)
}
