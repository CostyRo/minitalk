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

	obj.SetOptional("ifTrue", "ifFalse", core.NewObject(nil, ""))
	obj.SetOptional("ifFalse", "ifTrue", core.NewObject(nil, ""))

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
	obj.Set("not", !value, ObjectConstructor)
	obj.Set("ifTrue", func(other core.Object) interface{} {
		if other.Class != "CodeBlock" {
			return nil
		}
		noArgsVal, _ := other.Get("no_arguments")
		if noArgsVal.(int64) != 0 {
			return errors.NewValueError("CodeBlock must have no arguments").Object
		}
		if !value {
			other_, ok := obj.GetOptional("ifTrue","ifFalse")
			other = *other_
			if !ok || other.Class != "CodeBlock" {
				return core.Object{}
			}
		}
		valFn, _ := other.Get("value")
		result := valFn.(func(...core.Object) interface{})()
		return result
	})
	obj.Set("ifFalse", func(other core.Object) interface{} {
		if other.Class != "CodeBlock" {
			return nil
		}
		noArgsVal, _ := other.Get("no_arguments")
		if noArgsVal.(int64) != 0 {

			return errors.NewValueError("CodeBlock must have no arguments").Object
		}
		if value {
			other_, ok := obj.GetOptional("ifFalse","ifTrue")
			other = *other_
			if !ok || other.Class != "CodeBlock" {
				return core.Object{}
			}
		}
		valFn, _ := other.Get("value")
		result := valFn.(func(...core.Object) interface{})()
		return result
	})
	iVal := int64(0)
	if value {
		iVal = 1
	}
	obj.Set("toInteger", iVal, ObjectConstructor)
	obj.Set("toFloat", float64(iVal), ObjectConstructor)
	obj.Set("toBool", value, ObjectConstructor)
	obj.Set("toSymbol", fmt.Sprintf("%t", value), SymbolConstructor)
	obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
	obj.Set("toString", fmt.Sprintf("%t", value), ObjectConstructor)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to Array").Object)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)

	return &BoolObject{*obj}
}
