package types

import (
	"minitalk/tokens"
	"minitalk/types/core"
	"minitalk/types/errors"
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
	case []byte:
		return &NewByteArrayObject(v).Object
	case []core.Object:
		return &NewArrayObject(v).Object
	case core.NotImplementedObject:
		return &errors.NewNotImplementedError().Object
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

func convertToByteArray(arr []core.Object) ([]byte, bool) {
	bytes := make([]byte, len(arr))
	for i, el := range arr {
		if b, ok := el.Self.(byte); ok {
			bytes[i] = b
		} else if i64, ok := el.Self.(int64); ok && i64 >= 0 && i64 <= 255 {
			bytes[i] = byte(i64)
		} else {
			return nil, false
		}
	}
	return bytes, true
}

func StringToTokenType(s string) tokens.TokenType {
	switch s {
	case "Self":
		return tokens.Self_
	case "Super":
		return tokens.Super
	case "Nil":
		return tokens.Nil
	case "True":
		return tokens.True
	case "False":
		return tokens.False
	case "LParen":
		return tokens.LParen
	case "RParen":
		return tokens.RParen
	case "LBracket":
		return tokens.LBracket
	case "RBracket":
		return tokens.RBracket
	case "Period":
		return tokens.Period
	case "Semicolon":
		return tokens.Semicolon
	case "Colon":
		return tokens.Colon
	case "Pipe":
		return tokens.Pipe
	case "Caret":
		return tokens.Caret
	case "Plus":
		return tokens.Plus
	case "Minus":
		return tokens.Minus
	case "Star":
		return tokens.Star
	case "Slash":
		return tokens.Slash
	case "Ampersand":
		return tokens.Ampersand
	case "LessThan":
		return tokens.LessThan
	case "GreaterThan":
		return tokens.GreaterThan
	case "LessThanEqual":
		return tokens.LessThanEqual
	case "GreaterThanEqual":
		return tokens.GreaterThanEqual
	case "DoubleEquals":
		return tokens.DoubleEquals
	case "Assignment":
		return tokens.Assignment
	case "Identifier":
		return tokens.Identifier
	case "Integer":
		return tokens.Integer
	case "Float":
		return tokens.Float
	case "RadixNumber":
		return tokens.RadixNumber
	case "String":
		return tokens.String
	case "Symbol":
		return tokens.Symbol
	case "Character":
		return tokens.Character
	case "Array":
		return tokens.Array
	case "ByteArray":
		return tokens.ByteArray
	case "Comment":
		return tokens.Comment
	case "Whitespace":
		return tokens.Whitespace
	case "Illegal":
		return tokens.Illegal
	case "Error":
		return tokens.Error
	default:
		return tokens.Illegal
	}
}
