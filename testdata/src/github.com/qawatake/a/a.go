package a

func Wrap_NamedReturnValue() (err error) {
	defer Wrap(&err, "x") // ok
	return nil
}

func Wrap_LocalVariableNotNamedReturnValue() error {
	var err error
	defer Wrap(&err, "x") // want "named"
	return nil
}

func Wrapper_Wrap_NamedReturnValue() (err error) {
	var w wrapper
	defer w.Wrap(&err, "x") // ok
	return nil
}

func Wrapper_Wrap_LocalVariableNotReturnValue() {
	var w wrapper
	var err error
	defer w.Wrap(&err, "x") // want "named"
}

func Wrap(errp *error, msg string) {}

type wrapper struct{}

func (w wrapper) Wrap(errp *error, msg string) {}
