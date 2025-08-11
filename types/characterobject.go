package types

import (
	"fmt"

	"minitalk/types/core"
	"minitalk/types/errors"
)

type CharacterObject struct {
	core.Object
}

func NewCharacterObject(value rune) *CharacterObject {
	obj := core.NewObject(value, "Character")

	obj.Set("lt", func(other core.Object) interface{} {
		if other.Class != "Character" {
			return nil
		}
		return NewBoolObject(value < other.Self.(rune)).Object
	})
	obj.Set("gt", func(other core.Object) interface{} {
		if other.Class != "Character" {
			return nil
		}
		return NewBoolObject(value > other.Self.(rune)).Object
	})
	obj.Set("le", func(other core.Object) interface{} {
		if other.Class != "Character" {
			return nil
		}
		return NewBoolObject(value <= other.Self.(rune)).Object
	})
	obj.Set("ge", func(other core.Object) interface{} {
		if other.Class != "Character" {
			return nil
		}
		return NewBoolObject(value >= other.Self.(rune)).Object
	})
	obj.Set("eq", func(other core.Object) interface{} {
		if other.Class != "Character" {
			return nil
		}
		return NewBoolObject(value == other.Self.(rune)).Object
	})
	obj.Set("toInteger", int64(value), ObjectConstructor)
	obj.Set("toFloat", float64(value), ObjectConstructor)
	obj.Set("toBool", value != 0, ObjectConstructor)
	obj.Set("toSymbol", fmt.Sprintf("%c", value), SymbolConstructor)
	obj.Set("toCharacter", value, ObjectConstructor)
	obj.Set("toString", fmt.Sprintf("%c", value), ObjectConstructor)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to Array").Object)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to ByteArray").Object)

	return &CharacterObject{*obj}
}
