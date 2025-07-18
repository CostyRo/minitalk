package errors

func NewValueError(msgs ...string) *Error {
	msg := "wrong value"
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	return NewErrorObject(msg, "ValueError")
}