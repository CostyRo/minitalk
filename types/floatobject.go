package types

import (
	"fmt"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type FloatObject struct {
	core.Object
}

func NewFloatObject(value float64) *FloatObject {
	obj := core.NewObject(value, "Float")

	obj.Set("plus", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewFloatObject(value + val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewFloatObject(value + float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("minus", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewFloatObject(value - val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewFloatObject(value - float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("mul", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewFloatObject(value * val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewFloatObject(value * float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("div", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				if val == 0.0 {
					return errors.NewZeroDivisionError().Object
				}
				return NewFloatObject(value / val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				if val == 0 {
					return errors.NewZeroDivisionError().Object
				}
				return NewFloatObject(value / float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("lt", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(value < val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value < float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("gt", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(value > val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value > float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("le", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(value <= val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value <= float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("ge", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(value >= val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value >= float64(val)).Object
			}
		}
		return nil
	})
	obj.Set("eq", func(other core.Object) interface{} {
		switch other.Class {
		case "Float":
			if val, ok := other.Self.(float64); ok {
				return NewBoolObject(value == val).Object
			}
		case "Integer":
			if val, ok := other.Self.(int64); ok {
				return NewBoolObject(value == float64(val)).Object
			}
		}
		return nil
	})
	if val, ok := obj.Self.(float64); ok {
		obj.Set("toInteger", int64(val), ObjectConstructor)
		obj.Set("toFloat", val, ObjectConstructor)
		obj.Set("toBool", val != 0, ObjectConstructor)
		obj.Set("toSymbol", errors.NewTypeError("Invalid conversion to Symbol").Object)
		if val < 0 || val > 0x10FFFF {
			obj.Set("toCharacter", errors.NewValueError("Value is not in valid Unicode range 0..0x10FFFF").Object, nil)
		} else {
			obj.Set("toCharacter", rune(val), ObjectConstructor)
		}
		obj.Set("toString", fmt.Sprintf("%.10f", val), ObjectConstructor)
	}

	return &FloatObject{*obj}
}
