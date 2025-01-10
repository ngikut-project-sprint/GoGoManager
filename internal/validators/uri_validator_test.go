package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

func TestValidateURI_Valid(t *testing.T) {
	uri1 := "http://127.0.0.1:8080"
	uri2 := "https://127.0.0.1:8080"

	err1 := validators.ValidateURI(uri1)
	err2 := validators.ValidateURI(uri2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestValidateURI_WrongFormat(t *testing.T) {
	uri := "random@example.com"
	err := validators.ValidateURI(uri)
	assert.Error(t, err)
}

func TestValidateURI_NoScheme(t *testing.T) {
	uri := "127.0.0.1:8080"
	err := validators.ValidateURI(uri)
	assert.Error(t, err)
}

func TestValidateURI_NoHost(t *testing.T) {
	uri := "https://"
	err := validators.ValidateURI(uri)
	assert.Error(t, err)
}
