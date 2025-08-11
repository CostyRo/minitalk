package types

import (
	"fmt"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type IntegerObject struct {
	core.Object
}

func NewIntegerObject(value int64) *IntegerObject {
	obj := core.NewObject(value, "Integer")

	obj.Set("plus", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewIntegerObject(value + val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewFloatObject(float64(value) + val).Object
			}
		}
		return nil
	})
	obj.Set("minus", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewIntegerObject(value - val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewFloatObject(float64(value) - val).Object
			}
		}
		return nil
	})
	obj.Set("mul", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewIntegerObject(value * val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewFloatObject(float64(value) * val).Object
			}
		}
		return nil
	})
	obj.Set("div", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				if val == 0 {
					return errors.NewZeroDivisionError().Object
				}
				return NewIntegerObject(value / val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				if val == 0.0 {
					return errors.NewZeroDivisionError().Object
				}
				return NewFloatObject(float64(value) / val).Object
			}
		}
		return nil
	})
	obj.Set("lt", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value < val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(float64(value) < val).Object
			}
		}
		return nil
	})
	obj.Set("gt", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value > val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(float64(value) > val).Object
			}
		}
		return nil
	})
	obj.Set("le", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value <= val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(float64(value) <= val).Object
			}
		}
		return nil
	})
	obj.Set("ge", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value >= val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(float64(value) >= val).Object
			}
		}
		return nil
	})
	obj.Set("eq", func(other core.Object) interface{} {
		switch other.Class {
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value == val).Object
			}
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(float64(value) == val).Object
			}
		}
		return nil
	})
	obj.Set("toInteger", value, ObjectConstructor)
	obj.Set("toFloat", float64(value), ObjectConstructor)
	obj.Set("toBool", value != 0, ObjectConstructor)
	obj.Set("toSymbol", fmt.Sprintf("%d", value), SymbolConstructor)
	if value < 0 || value > 0x10FFFF {
		obj.Set("toCharacter", errors.NewValueError("Value is not in valid Unicode range 0..0x10FFFF").Object)
	} else {
		obj.Set("toCharacter", rune(value), ObjectConstructor)
	}
	obj.Set("toString", fmt.Sprintf("%d", value), ObjectConstructor)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to Array").Object)

	return &IntegerObject{*obj}
}
