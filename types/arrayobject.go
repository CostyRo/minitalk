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

	obj.SetOptional("at", "put", core.NewObject(nil, ""))

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
		if other.Class != "Integer" {
			return nil
		}
		idx, ok := other.Self.(int64)
		if !ok || idx < 0 || idx >= int64(len(elements)) {
			return errors.NewValueError(fmt.Sprintf("Index %d out of range", idx)).Object
		}
		putVal, _ := obj.GetOptional("at", "put")
		if putVal.Class != "" {
			elements[idx] = *putVal
		}
		return elements[idx]
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

		valFnVal, ok := other.Get("value")
		if !ok {
			return errors.NewValueError("CodeBlock missing value attribute").Object
		}
		callable, ok := valFnVal.(func(...core.Object) interface{})
		if !ok {
			return errors.NewValueError("Invalid code block value").Object
		}

		elements := obj.Self.([]core.Object)

		for i, elem := range elements {
			if noArgs == 1 {
				callable(elem)
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
				callable2(elem)
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
	bArr, ok := convertToByteArray(elements)
	if ok {
		obj.Set("toByteArray", bArr, ObjectConstructor)
	} else {
		obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)
	}
	obj.Set("toArray", elements, ObjectConstructor)

	return &ArrayObject{*obj}
}
