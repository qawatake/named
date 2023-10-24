package example

import (
	"errors"
	"fmt"
)

func Good() (err error) {
	err = DoSomething()
	defer Wrap(&err, "wrapped!!") // ok because a named return value is passed.
	return fmt.Errorf("error from Good: %w", err)
}

func Bad() error {
	err := DoSomething()
	defer Wrap(&err, "wrapped!!") // <- err is not a named return value.
	return fmt.Errorf("error from Bad: %w", err)
}

func Wrap(errp *error, msg string) {
	if *errp == nil {
		return
	}
	*errp = fmt.Errorf("%s: %w", msg, *errp)
}

func Example() {
	err := Good()
	fmt.Println(err)
	err = Bad()
	fmt.Println(err)
	// Output:
	// wrapped!!: error from Good: original error
	// error from Bad: original error
}

func DoSomething() error {
	return errors.New("original error")
}
