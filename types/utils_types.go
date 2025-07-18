package types

import (
	"minitalk/types/core"
)

func ObjectConstructor(val interface{}) *core.Object {
	switch v := val.(type) {
	case int64:
		return &NewIntegerObject(v).Object
	case float64:
		return &NewFloatObject(v).Object
	case bool:
		return &NewBoolObject(v).Object
	case string:
		return &NewStringObject(v).Object
	case rune:
		return &NewCharacterObject(v).Object
	default:
		return nil
	}
}

func SymbolConstructor(val interface{}) *core.Object {
	if s, ok := val.(string); ok {
		return &NewSymbolObject(s).Object
	}
	return nil
}