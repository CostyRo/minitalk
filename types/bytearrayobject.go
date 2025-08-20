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

	obj.SetOptional("at", "put", core.NewObject(nil, ""))

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
		if other.Class != "Integer" {
			return nil
		}
		idx, ok := other.Self.(int64)
		if !ok || idx < 0 || idx >= int64(len(data)) {
			return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
		}
		putVal, _ := obj.GetOptional("at", "put")
		if putVal.Class != "" {
			if putVal.Class != "Integer" {
				return errors.NewTypeError("ByteArray accepts only Integer values").Object
			}
			byteVal, ok := putVal.Self.(int64)
			if !ok || byteVal < 0 || byteVal > 255 {
				return errors.NewValueError(fmt.Sprintf("Invalid byte value: %v", putVal.Self)).Object
			}
			data[idx] = byte(byteVal)
			return NewIntegerObject(byteVal).Object
		}

		return NewIntegerObject(int64(data[idx])).Object
	})
	obj.Set("do", func(other core.Object) interface{} {
		if other.Class != "CodeBlock" {
			return nil
		}
		noArgsVal, ok := other.Get("no_arguments")
		if !ok {
			return errors.NewValueError("CodeBlock missing no_arguments attribute").Object
		}
		noArgs, ok := noArgsVal.(int64)
		if !ok || (noArgs != 1 && noArgs != 2) {
			return errors.NewValueError("CodeBlock must have 1 or 2 arguments").Object
		}

		valFn, _ := other.Get("value")
		callable, ok := valFn.(func(...core.Object) interface{})
		if !ok {
			return errors.NewValueError("Invalid code block value").Object
		}

		for i, b := range data {
			if noArgs == 1 {
				callable(NewIntegerObject(int64(b)).Object)
			} else {
				result := callable(NewIntegerObject(int64(i)).Object)
				nextBlock, ok := result.(core.Object)
				if !ok || nextBlock.Class != "CodeBlock" {
					continue
				}
				valFn2, _ := nextBlock.Get("value")
				callable2, ok := valFn2.(func(...core.Object) interface{})
				if !ok {
					continue
				}
				callable2(NewIntegerObject(int64(b)).Object)
			}
		}

		returnObj := NewBoolObject(true).Object
		returnObj.Set("!printable", false)
		return returnObj
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
