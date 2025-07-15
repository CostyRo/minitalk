package main

import (
	"regexp"
)

type TokenType int

const (
	Self_ TokenType = iota
	Super
	Nil
	True
	False
	LParen
	RParen
	LBracket
	RBracket
	Period
	Semicolon
	Colon
	Pipe
	Caret
	Plus
	Minus
	Star
	Slash
	Ampersand
	LessThan
	GreaterThan
	LessThanEqual
	GreaterThanEqual
	DoubleEquals
	Assignment
	Identifier
	Integer
	Float
	RadixNumber
	String
	Symbol
	Character
	Array
	ByteArray
	Comment
	Whitespace
	Illegal
	Error
)

type Token struct {
	Type  TokenType
	Value string
	Start int
	End   int
}

type tokenExpr struct {
	typ TokenType
	re  *regexp.Regexp
}

var tokenExprs = []tokenExpr{
	// Character (e.g. $x)
	{Character, regexp.MustCompile(`^\$.`)},

	// Symbol (e.g. #'symbol' or #identifier)
	{Symbol, regexp.MustCompile(`^#'([^']|'')*'|^#[a-zA-Z_][a-zA-Z0-9_]*`)},

	// Float (e.g. 123.45, 4.5e+6)
	{Float, regexp.MustCompile(`^(?:[0-9]+\.[0-9]+(?:[eE][+-]?[0-9]+)?|[0-9]+(?:[eE][+-]?[0-9]+))`)},

	// Radix number (e.g. 2r1010, 16rA000)
	{RadixNumber, regexp.MustCompile(`^[0-9]+r[0-9A-Fa-f]+`)},

	// Integer
	{Integer, regexp.MustCompile(`^[0-9]+`)},

	// Keywords
	{Self_, regexp.MustCompile(`^self\b`)},
	{Super, regexp.MustCompile(`^super\b`)},
	{Nil, regexp.MustCompile(`^nil\b`)},
	{True, regexp.MustCompile(`^true\b`)},
	{False, regexp.MustCompile(`^false\b`)},

	// Operators
	{LessThanEqual, regexp.MustCompile(`^<=`)},
	{GreaterThanEqual, regexp.MustCompile(`^>=`)},
	{DoubleEquals, regexp.MustCompile(`^==`)},
	{Assignment, regexp.MustCompile(`^:=`)},
	{LessThan, regexp.MustCompile(`^<`)},
	{GreaterThan, regexp.MustCompile(`^>`)},
	{Plus, regexp.MustCompile(`^\+`)},
	{Minus, regexp.MustCompile(`^-`)},
	{Star, regexp.MustCompile(`^\*`)},
	{Slash, regexp.MustCompile(`^/`)},
	{Ampersand, regexp.MustCompile(`^&`)},

	// Symbols
	{LParen, regexp.MustCompile(`^\(`)},
	{RParen, regexp.MustCompile(`^\)`)},
	{LBracket, regexp.MustCompile(`^\[`)},
	{RBracket, regexp.MustCompile(`^\]`)},
	{Period, regexp.MustCompile(`^\.`)},
	{Semicolon, regexp.MustCompile(`^;`)},
	{Colon, regexp.MustCompile(`^:`)},
	{Pipe, regexp.MustCompile(`^\|`)},
	{Caret, regexp.MustCompile(`^\^`)},

	// Identifier
	{Identifier, regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*`)},

	// String (e.g. 'hello', 'it''s fine')
	{String, regexp.MustCompile(`^'([^']|'')*'`)},

	// Array (e.g. #(...))
	{Array, regexp.MustCompile(`^#\([^)]*\)`)},

	// ByteArray (e.g. #[...])
	{ByteArray, regexp.MustCompile(`^#\[[^\]]*\]`)},

	// Comment (e.g. "...")
	{Comment, regexp.MustCompile(`^"[^"]*"`)},

	// Whitespace
	{Whitespace, regexp.MustCompile(`^[ \t\r\n]+`)},
}

func Lex(input string) []Token {
	var tokens []Token
	pos := 0

	for len(input) > 0 {
		matched := false

		for _, te := range tokenExprs {
			if loc := te.re.FindStringIndex(input); loc != nil && loc[0] == 0 {
				val := input[:loc[1]]
				tokens = append(tokens, Token{
					Type:  te.typ,
					Value: val,
					Start: pos,
					End:   pos + loc[1],
				})
				input = input[loc[1]:]
				pos += loc[1]
				matched = true
				break
			}
		}

		if !matched {
			tokens = append(tokens, Token{
				Type:  Illegal,
				Value: input[:1],
				Start: pos,
				End:   pos + 1,
			})
			input = input[1:]
			pos++
		}
	}

	return tokens
}
