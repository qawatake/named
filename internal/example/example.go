package example

import (
	"fmt"
)

func Good() (err error) {
	defer Wrap(&err, "x") // ok because a named return value is passed.
	return nil
}

func Bad() error {
	err := fmt.Errorf("x")
	defer Wrap(&err, "x") // <- err is not a named return value.
	return err
}

func Wrap(errp *error, msg string) {
	if errp == nil || *errp == nil {
		return
	}
	*errp = fmt.Errorf("%s: %w", msg, *errp)
}
