package validators

import (
	"fmt"
)

func ValidatePassword(password string, min int, max int) error {
	passwordLength := len(password)
	if passwordLength < min {
		return fmt.Errorf("Password too short (> %d characters)", min)
	}

	if passwordLength > max {
		return fmt.Errorf("Password too long (< %d characters)", max)
	}

	return nil
}
