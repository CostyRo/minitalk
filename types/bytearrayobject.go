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

	return &ByteArrayObject{*obj}
}
