package errors

func NewNameError(msgs ...string) *Error {
	msg := "name not defined"
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	return NewErrorObject(msg, "NameError")
}