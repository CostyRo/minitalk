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
	obj.Set("onError", func(other core.Object) interface{} {
		if other.Class != "CodeBlock" {
			return nil
		}

		noArgsVal, ok := other.Get("no_arguments")
		if !ok {
			return nil
		}
		noArgs, ok := noArgsVal.(int64)
		if !ok {
			return nil
		}

		valFnVal, ok := other.Get("value")
		if !ok {
			return nil
		}
		callable, ok := valFnVal.(func(...core.Object) interface{})
		if !ok {
			return nil
		}

		if noArgs == 0 {
			return callable()
		} else {
			return nil
		}
	})
	return &Error{*obj}
}
