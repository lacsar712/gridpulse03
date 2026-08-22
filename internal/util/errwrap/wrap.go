package errwrap

import (
	"errors"
	"fmt"
)

var ErrDenied = errors.New("denied")

func WrapDenied(op string) error {
	return fmt.Errorf("%s: %w", op, ErrDenied)
}

func Wrapf(err error, op string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}

func IsDenied(err error) bool {
	return errors.Is(err, ErrDenied)
}
