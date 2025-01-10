package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

func TestValidatePassword_Valid(t *testing.T) {
	password := "nanonano"
	err := validators.ValidatePassword(password, 4, 10)
	assert.NoError(t, err)
}

func TestValidatePassword_PasswordTooShort(t *testing.T) {
	password := "nan"
	err := validators.ValidatePassword(password, 4, 10)
	assert.Error(t, err)
}

func TestValidatePassword_PasswordTooLong(t *testing.T) {
	password := "nanonanonano"
	err := validators.ValidatePassword(password, 4, 10)
	assert.Error(t, err)
}
