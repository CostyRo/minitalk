package errors

func NewZeroDivisionError(msgs ...string) *Error {
	msg := "division by zero"
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	return NewErrorObject(msg, "ZeroDivisionError")
}