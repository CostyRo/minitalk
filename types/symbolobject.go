package types

import (
	"fmt"
	"strconv"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type SymbolObject struct {
	core.Object
}

func NewSymbolObject(name string) *SymbolObject {
	obj := core.NewObject(name, "Symbol")

	obj.Set("eq", func(other core.Object) interface{} {
		if other.Class != "Symbol" {
			return NewBoolObject(false).Object
		}
		return NewBoolObject(obj.String() == other.String()).Object
	})

	if iVal, err := strconv.ParseInt(name, 10, 64); err == nil {
		obj.Set("toInteger", iVal, ObjectConstructor)
		obj.Set("toFloat", float64(iVal), ObjectConstructor)
	} else {
		obj.Set("toInteger", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Integer", name)).Object)
		obj.Set("toFloat", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Float", name)).Object)
	}
	if bVal, err := strconv.ParseBool(name); err == nil {
		obj.Set("toBool", bVal, ObjectConstructor)
	} else {
		obj.Set("toBool", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Bool", name)).Object)
	}
	if len(name) == 1 {
		obj.Set("toCharacter", rune(name[0]), ObjectConstructor)
	} else {
		obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
	}
	str := name
	if len(name) >= 2 && name[0] == '\'' && name[len(name)-1] == '\'' {
		str = name[1:len(name)-1]
	}
	obj.Set("toString", str, ObjectConstructor)
	obj.Set("toSymbol", name, SymbolConstructor)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to Array").Object)

	return &SymbolObject{*obj}
}
