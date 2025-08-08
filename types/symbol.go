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
	if value, ok := obj.Self.(string); ok {
		symbolBody := value[1:]
		if len(symbolBody) >= 2 && symbolBody[0] == '\'' && symbolBody[len(symbolBody)-1] == '\'' {
			symbolBody = symbolBody[1 : len(symbolBody)-1]
		}
		if iVal, err := strconv.ParseInt(symbolBody, 10, 64); err == nil {
			obj.Set("toInteger", iVal, ObjectConstructor)
			obj.Set("toFloat", float64(iVal), ObjectConstructor)
		} else {
			obj.Set("toInteger", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Integer", symbolBody)).Object)
			obj.Set("toFloat", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Float", symbolBody)).Object)
		}
		if bVal, err := strconv.ParseBool(symbolBody); err == nil {
			obj.Set("toBool", bVal, ObjectConstructor)
		} else {
			obj.Set("toBool", errors.NewValueError(fmt.Sprintf("Cannot convert %s to Bool", symbolBody)).Object)
		}
		if len(symbolBody) == 1 {
			obj.Set("toCharacter", rune(symbolBody[0]), ObjectConstructor)
		} else {
			obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to Character").Object)
		}
		obj.Set("toString", symbolBody, ObjectConstructor)
		obj.Set("toSymbol", value, SymbolConstructor)
	}

	return &SymbolObject{*obj}
}
