package types

import (
	"fmt"

	"minitalk/types/core"
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
	if val, ok := obj.Self.(rune); ok {
		obj.Set("toInteger", int64(val), ObjectConstructor)
		obj.Set("toFloat", float64(val), ObjectConstructor)
		obj.Set("toBool", value != 0, ObjectConstructor)
		obj.Set("toSymbol", fmt.Sprintf("#%c", value), SymbolConstructor)
		obj.Set("toCharacter", val, ObjectConstructor)
		obj.Set("toString", fmt.Sprintf("%c", value), ObjectConstructor)
	}

	return &CharacterObject{*obj}
}
