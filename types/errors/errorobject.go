package errors

import "minitalk/types/core"

type Error struct {
	core.Object
}

func NewErrorObject(msg string, class ...string) *Error {
	className := "Error"
	if len(class) > 0 && class[0] != "" {
		className = class[0]
	}
	obj := core.NewObject(msg, className)
	return &Error{*obj}
}
