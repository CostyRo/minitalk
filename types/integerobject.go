package types

import (
	"minitalk/types/core"
	"minitalk/types/errors"
)

type IntegerObject struct {
	core.Object
}

func NewIntegerObject(value int64) *IntegerObject {
	obj := core.NewObject(value, "Integer")

	obj.Set("add", func(other interface{}) interface{} {
		switch b := other.(type) {
		case int64:
			return NewIntegerObject(value + b).Object
		case float64:
			return NewFloatObject(float64(value) + b).Object
		default:
			panic("Unsupported operand for IntegerObject.add")
		}
	})

	obj.Set("sub", func(other interface{}) interface{} {
		switch b := other.(type) {
		case int64:
			return NewIntegerObject(value - b).Object
		case float64:
			return NewFloatObject(float64(value) - b).Object
		default:
			panic("Unsupported operand for IntegerObject.sub")
		}
	})

	obj.Set("mul", func(other interface{}) interface{} {
		switch b := other.(type) {
		case int64:
			return NewIntegerObject(value * b).Object
		case float64:
			return NewFloatObject(float64(value) * b).Object
		default:
			panic("Unsupported operand for IntegerObject.mul")
		}
	})

	obj.Set("div", func(other interface{}) interface{} {
		switch b := other.(type) {
		case int64:
			if b == 0 {
				return errors.NewZeroDivisionError().Object
			}
			return NewIntegerObject(value / b).Object
		case float64:
			if b == 0.0 {
				return errors.NewZeroDivisionError().Object
			}
			return NewFloatObject(float64(value) / b).Object
		default:
			panic("Unsupported operand for IntegerObject.div")
		}
	})

	return &IntegerObject{*obj}
}
