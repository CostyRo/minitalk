package types

import "minitalk/types/core"

type BoolObject struct {
	core.Object
}

func NewBoolObject(value bool) *BoolObject {
	obj := core.NewObject(value, "Bool")
	return &BoolObject{*obj}
}
