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
	obj.Set("toInteger", int64(value), ObjectConstructor)
	obj.Set("toFloat", value, ObjectConstructor)
	obj.Set("toBool", value != 0, ObjectConstructor)
	obj.Set("toSymbol", errors.NewTypeError("Invalid conversion to Symbol").Object)
	if value < 0 || value > 0x10FFFF {
		obj.Set("toCharacter", errors.NewValueError("Value is not in valid Unicode range 0..0x10FFFF").Object)
	} else {
		obj.Set("toCharacter", rune(value), ObjectConstructor)
	}
	obj.Set("toString", fmt.Sprintf("%.10f", value), ObjectConstructor)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to Array").Object)

	return &FloatObject{*obj}
}
