package types

import (
	"minitalk/types/core"
	"minitalk/types/errors"
)

type FloatObject struct {
	core.Object
}

func NewFloatObject(value float64) *FloatObject {
	obj := core.NewObject(value, "Float")

	obj.Set("add", func(other core.Object) interface{} {
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

	obj.Set("sub", func(other core.Object) interface{} {
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

	return &FloatObject{*obj}
}
