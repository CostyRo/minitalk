package types

type FloatObject struct {
	Object
}

func NewFloatObject(value float64) *FloatObject {
	obj := NewObject(value, "float")

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

	return &FloatObject{*obj}
}
