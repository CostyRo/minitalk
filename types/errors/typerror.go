package errors

func NewTypeError(msgs ...string) *Error {
	msg := "type mismatch error"
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	return NewErrorObject(msg, "TypeError")
}