package errors

import "github.com/pkg/errors"

func New(message string) error {
	return errors.New(message)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
