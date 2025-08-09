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

	return &ArrayObject{*obj}
}
