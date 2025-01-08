package utils

import (
	"errors"
	"net/url"
)

func ValidateURI(uri string) error {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return errors.New("Invalid URI structure")
	}

	// Check if the scheme is allowed
	if parsedURI.Scheme != "http" && parsedURI.Scheme != "https" {
		return errors.New("Unsupported scheme (only http and https are allowed)")
	}

	// Check if the host is present
	if parsedURI.Host == "" {
		return errors.New("Missing host in URI")
	}

	return nil
}
