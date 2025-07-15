package errors

import "minitalk/types"

type Error struct {
	types.Object
}

func NewErrorObject(msg string) *Error {
	obj := types.NewObject(msg, "error")
	return &Error{*obj}
}
