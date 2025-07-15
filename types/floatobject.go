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

	obj.Set("add", func(other interface{}) interface{} {
		switch b := other.(type) {
		case float64:
			return NewFloatObject(value + b).Object
		case int64:
			return NewFloatObject(value + float64(b)).Object
		default:
			panic("Unsupported operand for FloatObject.add")
		}
	})

	obj.Set("sub", func(other interface{}) interface{} {
		switch b := other.(type) {
		case float64:
			return NewFloatObject(value - b).Object
		case int64:
			return NewFloatObject(value - float64(b)).Object
		default:
			panic("Unsupported operand for FloatObject.sub")
		}
	})

	obj.Set("mul", func(other interface{}) interface{} {
		switch b := other.(type) {
		case float64:
			return NewFloatObject(value * b).Object
		case int64:
			return NewFloatObject(value * float64(b)).Object
		default:
			panic("Unsupported operand for FloatObject.mul")
		}
	})

	obj.Set("div", func(other interface{}) interface{} {
		switch b := other.(type) {
		case float64:
			if b == 0.0 {
				return errors.NewZeroDivisionError().Object
			}
			return NewFloatObject(value / b).Object
		case int64:
			if b == 0 {
				return errors.NewZeroDivisionError().Object
			}
			return NewFloatObject(value / float64(b)).Object
		default:
			panic("Unsupported operand for FloatObject.div")
		}
	})

	return &FloatObject{*obj}
}
