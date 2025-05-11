package schemas

import "errors"

var (
	ErrorIncorrectExpression = errors.New("incorrect expression")
	ErrorDivisionByZero      = errors.New("division by zero")
	ErrorUnknownOperation    = errors.New("unknown operation")

	ErrorNotFound       = errors.New("not found")
	ErrorInternalServer = errors.New("internal server error")
)
