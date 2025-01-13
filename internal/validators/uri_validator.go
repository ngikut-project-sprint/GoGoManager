package validators

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
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

	components := strings.Split(uri, ".")
	if len(components) < 2 {
		return errors.New("Missing domain in URI")
	}

	emailRegex := `^[a-zA-Z0-9.-]+\.[a-zA-Z]{1,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(parsedURI.Host) {
		return errors.New("Invalid uri format")
	}

	return nil
}
