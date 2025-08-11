package main

import (
    "fmt"
    "runtime"

	"minitalk/tokens"
)

func Log(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("%s: at %s:%d\n", msg, file, line)
	}
}

func TokenTypeToString(t tokens.TokenType) string {
	switch t {
	case tokens.Self_:
		return "Self"
	case tokens.Super:
		return "Super"
	case tokens.Nil:
		return "Nil"
	case tokens.True:
		return "True"
	case tokens.False:
		return "False"
	case tokens.LParen:
		return "LParen"
	case tokens.RParen:
		return "RParen"
	case tokens.LBracket:
		return "LBracket"
	case tokens.RBracket:
		return "RBracket"
	case tokens.Period:
		return "Period"
	case tokens.Semicolon:
		return "Semicolon"
	case tokens.Colon:
		return "Colon"
	case tokens.Pipe:
		return "Pipe"
	case tokens.Caret:
		return "Caret"
	case tokens.Plus:
		return "Plus"
	case tokens.Minus:
		return "Minus"
	case tokens.Star:
		return "Star"
	case tokens.Slash:
		return "Slash"
	case tokens.Ampersand:
		return "Ampersand"
	case tokens.LessThan:
		return "LessThan"
	case tokens.GreaterThan:
		return "GreaterThan"
	case tokens.LessThanEqual:
		return "LessThanEqual"
	case tokens.GreaterThanEqual:
		return "GreaterThanEqual"
	case tokens.DoubleEquals:
		return "DoubleEquals"
	case tokens.Assignment:
		return "Assignment"
	case tokens.Identifier:
		return "Identifier"
	case tokens.Integer:
		return "Integer"
	case tokens.Float:
		return "Float"
	case tokens.RadixNumber:
		return "RadixNumber"
	case tokens.String:
		return "String"
	case tokens.Symbol:
		return "Symbol"
	case tokens.Character:
		return "Character"
	case tokens.Array:
		return "Array"
	case tokens.ByteArray:
		return "ByteArray"
	case tokens.Comment:
		return "Comment"
	case tokens.Whitespace:
		return "Whitespace"
	case tokens.Illegal:
		return "Illegal"
	case tokens.Error:
		return "Error"
	default:
		return "Unknown"
	}
}
