package a

import (
	"a/b"
	"errors"
)

func Wrap_NamedReturnValue() (err error) {
	defer Wrap(&err, "x") // ok
	return nil
}

func Wrap_LocalVariableDeclared() error {
	var err error
	defer Wrap(&err, "x") // want "named"
	return nil
}

func Wrap_LocalVariableAssign() error {
	err := DoSomething()
	defer Wrap(&err, "x") // want "named"
	return err
}

func Wrap_FieldOfNamedReturnValue() (s S) {
	defer Wrap(&s.Field, "x") // ok
	return
}

func Wrap_FieldOfLocalVariable() S {
	var s S
	defer Wrap(&s.Field, "x") // want "named"
	return s
}

func Wrap_NillLiteral() (err error) {
	defer Wrap(nil, "x") // want "named"
	return nil
}

func B_Wrap_MethodNamedReturnValue() (err error) {
	defer b.Wrap(&err, "x") // ok
	return nil
}

func B_Wrap_MethodLocalVariableDeclared() error {
	var err error
	defer b.Wrap(&err, "x") // want "named"
	return nil
}

func B_Wrap_MethodNilLiteral() (err error) {
	defer b.Wrap(nil, "x") // want "named"
	return nil
}

func WrapAny_NamedReturnValue() (err error) {
	defer WrapAny(&err) // ok
	return nil
}

func WrapAny_UnaryNonPointer() (err error) {
	defer WrapAny(+1) // want "named"
	return nil
}

func WrapAny_UnaryPointerStar() (err error) {
	var x *int
	defer WrapAny(&*x) // want "named"
	return nil
}

func Wrap_Anonymous() {
	func() (err error) {
		defer Wrap(&err, "x") // ok
		return
	}()
}

func Wrap_NotInnerMost() (err error) {
	func() {
		defer Wrap(&err, "x") // want "named"
		return
	}()
	return
}

func Wrap_AnonymousDoubleNamedReturnValues() (err error) {
	func() (err error) {
		defer Wrap(&err, "x") // ok
		return
	}()
	return
}

func B_Wrap_NotInnerMost() (err error) {
	func() {
		defer b.Wrap(&err, "x") // want "named"
		return
	}()
	return
}

func Wrapper_Wrap_NamedReturnValue() (err error) {
	var w wrapper
	defer w.Wrap(&err, "x") // ok
	return nil
}

func Wrapper_Wrap_LocalVariableDeclared() {
	var w wrapper
	var err error
	defer w.Wrap(&err, "x") // want "named"
}

func Wrap_LocalVariableWithTheSameNameInFor() (err error) {
	for i := 0; i < 10; i++ {
		err := errors.New("x")
		defer Wrap(&err, "x") // want "named"
	}
	return
}

func Wrap(errp *error, msg string) {}

func WrapAny(v any) {}

type wrapper struct{}

func (w wrapper) Wrap(errp *error, msg string) {}

func DoSomething() error {
	return errors.New("error from DoSomething")
}

type S struct {
	Field error
}
