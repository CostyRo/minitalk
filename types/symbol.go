package types

import "minitalk/types/core"

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

	return &SymbolObject{*obj}
}
