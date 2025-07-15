package types

type IntegerObject struct {
	Object
}

func NewIntegerObject(value int64) *IntegerObject {
	obj := NewObject(value, "integer")

	obj.Set("add", func(other interface{}) interface{} {
		if b, ok := other.(int64); ok {
			return NewIntegerObject(value + b).Object
		}
		if b, ok := other.(float64); ok {
			return NewFloatObject(float64(value) + b).Object
		}
		panic("Unsupported operand for IntegerObject.add")
	})

	obj.Set("sub", func(other interface{}) interface{} {
		if b, ok := other.(int64); ok {
			return NewIntegerObject(value - b).Object
		}
		if b, ok := other.(float64); ok {
			return NewFloatObject(float64(value) - b).Object
		}
		panic("Unsupported operand for IntegerObject.sub")
	})

	obj.Set("mul", func(other interface{}) interface{} {
		if b, ok := other.(int64); ok {
			return NewIntegerObject(value * b).Object
		}
		if b, ok := other.(float64); ok {
			return NewFloatObject(float64(value) * b).Object
		}
		panic("Unsupported operand for IntegerObject.mul")
	})

	return &IntegerObject{*obj}
}
