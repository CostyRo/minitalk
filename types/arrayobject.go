package types

import (
	"fmt"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type ArrayObject struct {
	core.Object
}

func NewArrayObject(elements []core.Object) *ArrayObject {
	obj := core.NewObject(elements, "Array")

	obj.Set("plus", func(other core.Object) interface{} {
		if other.Class == "Array" {
			if val, ok := other.Self.([]core.Object); ok {
				newData := append(append([]core.Object{}, elements...), val...)
				return NewArrayObject(newData).Object
			}
		}
		return nil
	})

	obj.Set("at", func(other core.Object) interface{} {
		if other.Class == "Integer" {
			if idx, ok := other.Self.(int64); ok {
				if idx < 0 || idx >= int64(len(elements)) {
					return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
				}
				return elements[idx]
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
	bArr, ok := convertToByteArray(elements)
	if ok {
		obj.Set("toByteArray", bArr, ObjectConstructor)
	} else {
		obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)
	}
	obj.Set("toArray", elements, ObjectConstructor)

	return &ArrayObject{*obj}
}
