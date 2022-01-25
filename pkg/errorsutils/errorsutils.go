package errorsutils

import "fmt"

// Wrap is a drop-in replacement for errors.Wrap (https://github.com/pkg/errors) using the std's fmt.Errorf().
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
