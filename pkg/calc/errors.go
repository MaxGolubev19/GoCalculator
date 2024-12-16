package calc

import "errors"

var (
	ErrorDivisionByZero      = errors.New("division by zero")
	ErrorIncorrectExpression = errors.New("incorrect expression")
)
