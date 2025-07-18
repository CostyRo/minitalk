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

	obj.Set("add", func(other core.Object) interface{} {
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
	obj.Set("sub", func(other core.Object) interface{} {
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
	if val, ok := obj.Self.(int64); ok {
		obj.Set("toInteger", val, ObjectConstructor)
		obj.Set("toFloat", float64(val), ObjectConstructor)
		obj.Set("toBool", val != 0, ObjectConstructor)
		obj.Set("toSymbol", fmt.Sprintf("#%d", val), SymbolConstructor)
		if val < 0 || val > 0x10FFFF {
			obj.Set("toCharacter", errors.NewValueError("Value is not in valid Unicode range 0..0x10FFFF").Object, nil)
		} else {
			obj.Set("toCharacter", rune(val), ObjectConstructor)
		}
		obj.Set("toString", fmt.Sprintf("%d", val), ObjectConstructor)
	}

	return &IntegerObject{*obj}
}
