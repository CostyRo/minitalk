package errors

func NewZeroDivisionError() *Error {
	return NewErrorObject("ZeroDivisionError: division by zero")
}