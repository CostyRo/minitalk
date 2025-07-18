package types

import (
	"fmt"

	"minitalk/types/core"
	"minitalk/types/errors"
)

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
	if _, ok := obj.Self.(bool); ok {
		iVal := int64(0)
		if value {
			iVal = 1
		}
		obj.Set("toInteger", iVal, ObjectConstructor)
		obj.Set("toFloat", float64(iVal), ObjectConstructor)
		obj.Set("toBool", value, ObjectConstructor)
		obj.Set("toSymbol", fmt.Sprintf("#%t", value), SymbolConstructor)
		obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
		obj.Set("toString", fmt.Sprintf("%t", value), ObjectConstructor)
	}

	return &BoolObject{*obj}
}
