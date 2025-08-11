package tokens

import (
	"testing"
)

func collectTokens(input string) []Token {
	var result []Token
	for _, tok := range Lex(input) {
		if tok.Type != Whitespace {
			result = append(result, Token{
				Type:  tok.Type,
				Value: tok.Value,
			})
		}
	}
	return result
}

func assertTokenSeq(t *testing.T, actual []Token, expected []Token) {
	if len(actual) != len(expected) {
		t.Fatalf("Expected %d tokens but got %d\nFound tokens: %v", len(expected), len(actual), actual)
	}
	for i, tok := range actual {
		if tok.Type != expected[i].Type {
			t.Errorf("Token at position %d: expected %v but got %v", i, expected[i].Type, tok.Type)
		}
		if tok.Value != expected[i].Value {
			t.Errorf("Text at position %d: expected '%s' but got '%s'", i, expected[i].Value, tok.Value)
		}
	}
}

func TestPostcard(t *testing.T) {
	input := `exampleWithNumber: x
    | y |
    true & false not & (nil isNil) ifFalse: [self halt].
    y := self size + super size.
    #($a #a 'a' 1 1.0)
        do: [ :each |
            Transcript show: (each class name);
                       show: ' '].
    ^x < y`

	expected := []Token{
		{Identifier, "exampleWithNumber", 0, 0},
		{Colon, ":", 0, 0},
		{Identifier, "x", 0, 0},
		{Pipe, "|", 0, 0},
		{Identifier, "y", 0, 0},
		{Pipe, "|", 0, 0},
		{True, "true", 0, 0},
		{Ampersand, "&", 0, 0},
		{False, "false", 0, 0},
		{Identifier, "not", 0, 0},
		{Ampersand, "&", 0, 0},
		{LParen, "(", 0, 0},
		{Nil, "nil", 0, 0},
		{Identifier, "isNil", 0, 0},
		{RParen, ")", 0, 0},
		{Identifier, "ifFalse", 0, 0},
		{Colon, ":", 0, 0},
		{LBracket, "[", 0, 0},
		{Self_, "self", 0, 0},
		{Identifier, "halt", 0, 0},
		{RBracket, "]", 0, 0},
		{Period, ".", 0, 0},
		{Identifier, "y", 0, 0},
		{Assignment, ":=", 0, 0},
		{Self_, "self", 0, 0},
		{Identifier, "size", 0, 0},
		{Plus, "+", 0, 0},
		{Super, "super", 0, 0},
		{Identifier, "size", 0, 0},
		{Period, ".", 0, 0},
		{Array, "#($a #a 'a' 1 1.0)", 0, 0},
		{Identifier, "do", 0, 0},
		{Colon, ":", 0, 0},
		{LBracket, "[", 0, 0},
		{Colon, ":", 0, 0},
		{Identifier, "each", 0, 0},
		{Pipe, "|", 0, 0},
		{Identifier, "Transcript", 0, 0},
		{Identifier, "show", 0, 0},
		{Colon, ":", 0, 0},
		{LParen, "(", 0, 0},
		{Identifier, "each", 0, 0},
		{Identifier, "class", 0, 0},
		{Identifier, "name", 0, 0},
		{RParen, ")", 0, 0},
		{Semicolon, ";", 0, 0},
		{Identifier, "show", 0, 0},
		{Colon, ":", 0, 0},
		{String, "' '", 0, 0},
		{RBracket, "]", 0, 0},
		{Period, ".", 0, 0},
		{Caret, "^", 0, 0},
		{Identifier, "x", 0, 0},
		{LessThan, "<", 0, 0},
		{Identifier, "y", 0, 0},
	}

	actual := collectTokens(input)
	assertTokenSeq(t, actual, expected)
}

func TestExtraCode(t *testing.T) {
	input := `
        < > <= >= == :=
        42 123.45 1.2e3 16rA000 2r1010
        'he''llo' #1 #'symbol' $x
        #($a #a 'b' 2 2.0 #(1)) #[1 2 3]
        "This is a comment"
    `

	expected := []Token{
		{LessThan, "<", 0, 0},
		{GreaterThan, ">", 0, 0},
		{LessThanEqual, "<=", 0, 0},
		{GreaterThanEqual, ">=", 0, 0},
		{DoubleEquals, "==", 0, 0},
		{Assignment, ":=", 0, 0},
		{Integer, "42", 0, 0},
		{Float, "123.45", 0, 0},
		{Float, "1.2e3", 0, 0},
		{RadixNumber, "16rA000", 0, 0},
		{RadixNumber, "2r1010", 0, 0},
		{String, "'he''llo'", 0, 0},
		{Symbol, "#1", 0, 0},
		{Symbol, "#'symbol'", 0, 0},
		{Character, "$x", 0, 0},
		{Array, "#($a #a 'b' 2 2.0 #(1))", 0, 0},
		{ByteArray, "#[1 2 3]", 0, 0},
		{Comment, "\"This is a comment\"", 0, 0},
	}

	actual := collectTokens(input)
	assertTokenSeq(t, actual, expected)
}
