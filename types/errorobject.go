package types

type ErrorObject struct {
	Object
}

func NewErrorObject(message string) *ErrorObject {
	obj := NewObject(message, "error")
	return &ErrorObject{*obj}
}

func (e *ErrorObject) GetMessage() string {
	msg, _ := GetSelfValue[string](&e.Object)
	return msg
}
