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
	obj.SetOptional("at", "insert", core.NewObject(nil, ""))

	obj.Set("plus", func(other core.Object) interface{} {
		if other.Class == "ByteArray" {
			if val, ok := other.Self.([]byte); ok {
				newData := append(append([]byte{}, data...), val...)
				return NewByteArrayObject(newData).Object
			}
		}
		return nil
	})
	obj.Set("eq", func(other core.Object) interface{} {
		if other.Class != "ByteArray" {
			return nil
		}
		val, ok := other.Self.([]byte)
		if !ok {
			return nil
		}
		if len(data) != len(val) {
			return NewBoolObject(false).Object
		}
		for i := range data {
			if data[i] != val[i] {
				return NewBoolObject(false).Object
			}
		}
		return NewBoolObject(true).Object
	})
	obj.Set("size", func() core.Object { return NewIntegerObject(int64(len(data))).Object })
	obj.Set("reversed", func() core.Object {
		data := obj.Self.([]byte)
		n := len(data)
		newData := make([]byte, n)
		for i, v := range data {
			newData[n-1-i] = v
		}
		return NewByteArrayObject(newData).Object
	})
	obj.Set("removeAt", func(other core.Object) interface{} {
		if other.Class != "Integer" {
			return nil
		}

		idx, ok := other.Self.(int64)
		data := obj.Self.([]byte)
		if !ok || idx < 0 || idx >= int64(len(data)) {
			return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
		}

		popped := data[idx]
		rest := make([]*core.Object, 0, len(data)-1)
		for i, b := range data {
			if int64(i) != idx {
				rest = append(rest, &NewIntegerObject(int64(b)).Object)
			}
		}

		return NewArrayObject([]*core.Object{
			&NewIntegerObject(int64(popped)).Object,
			&NewArrayObject(rest).Object,
		}).Object
	})
	obj.Set("at", func(other core.Object) interface{} {
		if other.Class != "Integer" {
			return nil
		}
		idx, ok := other.Self.(int64)
		if !ok || idx < 0 || idx > int64(len(data)) {
			return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
		}

		insertVal, _ := obj.GetOptional("at", "insert")
		if insertVal.Class != "" {
			if insertVal.Class != "Integer" {
				return nil
			}
			byteVal, ok := insertVal.Self.(int64)
			if !ok || byteVal < 0 || byteVal > 255 {
				return errors.NewValueError(fmt.Sprintf("Invalid byte value: %v", insertVal.Self)).Object
			}
			if idx == int64(len(data)) {
				data = append(data, byte(byteVal))
			} else {
				newData := make([]byte, len(data)+1)
				copy(newData[:idx], data[:idx])
				newData[idx] = byte(byteVal)
				copy(newData[idx+1:], data[idx:])
				data = newData
			}
			return NewByteArrayObject(data).Object
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
			if idx >= int64(len(data)) {
				return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
			}
			data[idx] = byte(byteVal)
			return NewByteArrayObject(data).Object
		}

		if idx >= int64(len(data)) {
			return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
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
	obj.Set("map", func(other core.Object) interface{} {
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
		valFnVal, ok := other.Get("value")
		if !ok {
			return errors.NewValueError("CodeBlock missing value attribute").Object
		}
		callable, ok := valFnVal.(func(...core.Object) interface{})
		if !ok {
			return errors.NewValueError("Invalid code block value").Object
		}

		data := obj.Self.([]byte)
		mapped := make([]*core.Object, len(data))

		for i, b := range data {
			elemObj := NewIntegerObject(int64(b)).Object
			if noArgs == 1 {
				res := callable(elemObj)
				if o, ok := res.(core.Object); ok {
					mapped[i] = &o
				} else {
					return nil
				}
			} else {
				result := callable(NewIntegerObject(int64(i)).Object)
				nextBlock, ok := result.(core.Object)
				if !ok || nextBlock.Class != "CodeBlock" {
					continue
				}
				valFn2Val, ok := nextBlock.Get("value")
				if !ok {
					continue
				}
				callable2, ok := valFn2Val.(func(...core.Object) interface{})
				if !ok {
					continue
				}
				res := callable2(elemObj)
				if o, ok := res.(core.Object); ok {
					mapped[i] = &o
				} else {
					return nil
				}
			}
		}
		return NewArrayObject(mapped).Object
	})
	obj.Set("toInteger", errors.NewTypeError("Invalid conversion to Integer").Object)
	obj.Set("toFloat", errors.NewTypeError("Invalid conversion to Float").Object)
	obj.Set("toBool", errors.NewTypeError("Invalid conversion to Bool").Object)
	obj.Set("toSymbol", errors.NewTypeError("Invalid conversion to Symbol").Object)
	obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
	obj.Set("toString", func() core.Object { return NewStringObject(obj.String()).Object })
	obj.Set("toByteArray", data, ObjectConstructor)
	elements := make([]*core.Object, len(data))
	for i, b := range data {
		elements[i] = &NewIntegerObject(int64(b)).Object
	}
	obj.Set("toArray", NewArrayObject(elements).Object)

	return &ByteArrayObject{*obj}
}
