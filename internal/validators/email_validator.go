package validators

import (
	"errors"
	"regexp"
	"strings"
)

type EmailValidator interface {
	ValidateEmail(email string) error
}

func ValidateEmail(email string) error {
	// Regex structure validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{1,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("Invalid email format")
	}

	// Parsing email using net/mail
	// _, err := mail.ParseAddress(email)
	// if err != nil {
	// 	return errors.New("Invalid email structure (RFC compliance)")
	// }

	// Domain validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("Invalid email format (missing domain)")
	}

	// domain := parts[1]
	// maxRecords, err := net.LookupMX(domain)
	// if err != nil || len(maxRecords) == 0 {
	// 	return errors.New("Domain does not exist or no MX records found")
	// }

	return nil
}
