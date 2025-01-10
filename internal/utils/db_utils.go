package utils

import (
	"fmt"

	"github.com/lib/pq"
)

func UniqueConstraintError(err error) error {
	pqErr, ok := err.(*pq.Error)
	if ok && pqErr.Code == "23505" {
		return fmt.Errorf("Unique constraint violated: %s", pqErr.Constraint)
	}
	return nil
}
