package a

import (
	"a/b"
	"fmt"
)

func f() (err error) {
	defer Wrap(&err, "x")
	fmt.Println("x")
	return nil
}

func g() error {
	var err error
	defer Wrap(&err, "x") // want "named"
	defer Wrap(nil, "x")  // want "named"
	defer Wrap(           // want "named"
		func() *error {
			var err error
			return &err
		}(),
		"x")
	return nil
}

func h() {
	defer Wrap(nil, "x") // want "named"
}

func f2() (err error) {
	defer b.Wrap(&err, "x")
	fmt.Println("x")
	return nil
}

func g2() error {
	var err error
	defer b.Wrap(&err, "x") // want "named"
	defer b.Wrap(nil, "x")  // want "named"
	defer b.Wrap(           // want "named"
		func() *error {
			var err error
			return &err
		}(),
		"x")
	return nil
}

func h2() {
	defer b.Wrap(nil, "x") // want "named"
}

func f3() (err error) {
	defer WrapAny(&err)
	return nil
}

func f4() error {
	var x *int
	defer WrapAny(+1)  // want "named"
	defer WrapAny(&*x) // want "named"
	return nil
}

func f5() {
	func() (err error) {
		defer Wrap(&err, "x") // ok
		return
	}()
}

func f6() (err error) {
	func() {
		defer Wrap(&err, "x") // want "named"
		return
	}()
	return
}

func f8() (err error) {
	func() (err error) {
		defer Wrap(&err, "x") // ok
		return
	}()
	return
}

func f7() (err error) {
	func() {
		defer b.Wrap(&err, "x") // want "named"
		return
	}()
	return
}

// todo method

func Wrap(errp *error, msg string) {}

func WrapAny(v any) {}
