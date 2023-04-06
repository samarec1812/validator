package validator

import "errors"

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	// TODO: implement this
	switch len(v) {
	case 0:
		return ""
	case 1:
		return v[0].Err.Error()
	}
	err := v[0].Err
	for _, e := range v[1:] {
		err = errors.Join(err, e.Err)
	}
	return err.Error()
}
