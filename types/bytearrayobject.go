package types

import (
	"fmt"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type ByteArrayObject struct {
	core.Object
}

func NewByteArrayObject(data []byte) *ByteArrayObject {
	obj := core.NewObject(data, "ByteArray")

	obj.Set("plus", func(other core.Object) interface{} {
		if other.Class == "ByteArray" {
			if val, ok := other.Self.([]byte); ok {
				newData := append(append([]byte{}, data...), val...)
				return NewByteArrayObject(newData).Object
			}
		}
		return nil
	})

	obj.Set("at", func(other core.Object) interface{} {
		if other.Class == "Integer" {
			if idx, ok := other.Self.(int64); ok {
				if idx < 0 || idx >= int64(len(data)) {
					return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
				}
				return NewIntegerObject(int64(data[idx])).Object
			}
		}
		return nil
	})

	
	obj.Set("toInteger", errors.NewTypeError("Invalid conversion to Integer").Object)
	obj.Set("toFloat", errors.NewTypeError("Invalid conversion to Float").Object)
	obj.Set("toBool", errors.NewTypeError("Invalid conversion to Bool").Object)
	obj.Set("toSymbol", errors.NewTypeError("Invalid conversion to Symbol").Object)
	obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
	obj.Set("toString", errors.NewTypeError("Invalid conversion to String").Object)
	obj.Set("toByteArray", data, ObjectConstructor)
	elements := make([]core.Object, len(data))
	for i, b := range data {
		elements[i] = NewIntegerObject(int64(b)).Object
	}
	obj.Set("toArray", NewArrayObject(elements).Object)

	return &ByteArrayObject{*obj}
}
