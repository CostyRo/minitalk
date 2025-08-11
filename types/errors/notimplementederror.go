package errors

func NewNotImplementedError(msgs ...string) *Error {
	msg := "not implemented"
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	return NewErrorObject(msg, "NotImplementedError")
}
