package errs

import (
	"errors"
	"fmt"
)

var (
	NotFoundError = errors.New("not found")
	ConflictError = errors.New("conflict")
)

func NewNotFoundError(msg string) error {
	return fmt.Errorf("%w: %s", NotFoundError, msg)
}

func NewConflictError(msg string) error {
	return fmt.Errorf("%w: %s", ConflictError, msg)
}
