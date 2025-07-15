package errors

func NewZeroDivisionError() *Error {
	return NewErrorObject("division by zero","ZeroDivisionError")
}