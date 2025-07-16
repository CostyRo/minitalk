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

	return &BoolObject{*obj}
}
