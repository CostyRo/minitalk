package errors

import "minitalk/types/core"

func NewValueError(msgs ...string) *Error {
	msg := "wrong value"
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	err := NewErrorObject(msg, "ValueError")
	err.Set("onValueError", func(other core.Object) interface{} {
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
	return err
}
