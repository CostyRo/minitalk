package types

import "minitalk/types/core"

type BoolObject struct {
	core.Object
}

func NewBoolObject(value bool) *BoolObject {
	obj := core.NewObject(value, "Bool")

	obj.Set("and", func(other core.Object) interface{} {
		if other.Class != "Bool" {
			return nil
		}
		otherVal, ok := other.Self.(bool)
		if !ok {
			return nil
		}
		return NewBoolObject(value && otherVal).Object
	})

	obj.Set("eq", func(other core.Object) interface{} {
		if other.Class != "Bool" {
			return NewBoolObject(false).Object
		}
		otherVal, ok := other.Self.(bool)
		if !ok {
			return NewBoolObject(false).Object
		}
		return NewBoolObject(value == otherVal).Object
	})

	return &BoolObject{*obj}
}
