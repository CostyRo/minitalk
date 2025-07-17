package types

import "minitalk/types/core"

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

	return &CharacterObject{*obj}
}
