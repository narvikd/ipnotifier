package errorsutils

import "fmt"

func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
