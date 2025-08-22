package types

import (
	"fmt"
	"strconv"
	"strings"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type StringObject struct {
	core.Object
}

func NewStringObject(value string) *StringObject {
	obj := core.NewObject(value, "String")

	obj.Set("plus", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewStringObject(value + other.Self.(string)).Object
	})
	obj.Set("lt", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value < other.Self.(string)).Object
	})
	obj.Set("gt", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value > other.Self.(string)).Object
	})
	obj.Set("le", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value <= other.Self.(string)).Object
	})
	obj.Set("ge", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value >= other.Self.(string)).Object
	})
	obj.Set("eq", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value == other.Self.(string)).Object
	})
	obj.Set("splitBy", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		sep := other.Self.(string)
		if sep == "" {
			sep = "\n"
		}
		parts := strings.Split(value, sep)
		arr := make([]*core.Object, len(parts))
		for i, p := range parts {
			arr[i] = &NewStringObject(strings.ReplaceAll(p, "\r", "")).Object
		}
		return NewArrayObject(arr).Object
	})
	if iVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		obj.Set("toInteger", iVal, ObjectConstructor)
	} else {
		obj.Set("toInteger", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Integer", value)).Object)
	}
	if fVal, err := strconv.ParseFloat(value, 64); err == nil {
		obj.Set("toFloat", fVal, ObjectConstructor)
	} else {
		obj.Set("toFloat", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Float", value)).Object)
	}
	if bVal, err := strconv.ParseBool(value); err == nil {
		obj.Set("toBool", bVal, ObjectConstructor)
	} else {
		obj.Set("toBool", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Bool", value)).Object)
	}
	if len(value) == 1 {
		obj.Set("toCharacter", rune(value[0]), ObjectConstructor)
	} else {
		obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
	}
	obj.Set("toString", value, ObjectConstructor)
	obj.Set("toSymbol", value, SymbolConstructor)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to Array").Object)

	return &StringObject{*obj}
}
