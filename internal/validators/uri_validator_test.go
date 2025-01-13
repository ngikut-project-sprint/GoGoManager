package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

func TestValidateURI_Valid(t *testing.T) {
	uri1 := "http://projectsprint-bucket-public-read.s3.ap-southeast-1.amazonaws.com/ngikut/gogomanager/c7feabe3-1bb7-4bd6-838d-3994dac2fa22.png"
	uri2 := "https://projectsprint-bucket-public-read.s3.ap-southeast-1.amazonaws.com/ngikut/gogomanager/c7feabe3-1bb7-4bd6-838d-3994dac2fa22.png"

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
