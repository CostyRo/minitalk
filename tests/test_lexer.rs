use smalltalk_compiler::lexer::{Lexer, Token};

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_postcard() {
        let postcard = 
"exampleWithNumber: x
    | y |
    true & false not & (nil isNil) ifFalse: [self halt].
    y := self size + super size.
    #($a #a 'a' 1 1.0)
        do: [ :each |
            Transcript show: (each class name);
                       show: ' '].
    ^x < y";

        let mut lexer = Lexer::new(postcard);
        let tokens: Vec<(Token, &str)> = collect_tokens(&mut lexer);

        let expected_sequence = vec![
            (Token::Identifier, "exampleWithNumber"),
            (Token::Colon, ":"),
            (Token::Identifier, "x"),
            (Token::Pipe, "|"),
            (Token::Identifier, "y"),
            (Token::Pipe, "|"),
            (Token::True, "true"),
            (Token::Ampersand, "&"),
            (Token::False, "false"),
            (Token::Identifier, "not"),
            (Token::Ampersand, "&"),
            (Token::LParen, "("),
            (Token::Nil, "nil"),
            (Token::Identifier, "isNil"),
            (Token::RParen, ")"),
            (Token::Identifier, "ifFalse"),
            (Token::Colon, ":"),
            (Token::LBracket, "["),
            (Token::Self_, "self"),
            (Token::Identifier, "halt"),
            (Token::RBracket, "]"),
            (Token::Period, "."),
            (Token::Identifier, "y"),
            (Token::Assignment, ":="),
            (Token::Self_, "self"),
            (Token::Identifier, "size"),
            (Token::Plus, "+"),
            (Token::Super, "super"),
            (Token::Identifier, "size"),
            (Token::Period, "."),
            (Token::Array, "#($a #a 'a' 1 1.0)"),
            (Token::Identifier, "do"),
            (Token::Colon, ":"),
            (Token::LBracket, "["),
            (Token::Colon, ":"),
            (Token::Identifier, "each"),
            (Token::Pipe, "|"),
            (Token::Identifier, "Transcript"),
            (Token::Identifier, "show"),
            (Token::Colon, ":"),
            (Token::LParen, "("),
            (Token::Identifier, "each"),
            (Token::Identifier, "class"),
            (Token::Identifier, "name"),
            (Token::RParen, ")"),
            (Token::Semicolon, ";"),
            (Token::Identifier, "show"),
            (Token::Colon, ":"),
            (Token::String, "' '"),
            (Token::RBracket, "]"),
            (Token::Period, "."),
            (Token::Caret, "^"),
            (Token::Identifier, "x"),
            (Token::LessThan, "<"),
            (Token::Identifier, "y"),
        ];

        assert_token_seq(&tokens, &expected_sequence);
    }

    #[test]
    fn test_extra_code() {
        let extra_code = r#"
        < > <= >= == :=
        -42 123.45 1.2e3 16rA000 2r1010
        'he''llo' #'symbol' $x
        #($a #a 'b' 2 2.0) #[1 2 3]
        "This is a comment"
        "#;

        let mut lexer = Lexer::new(extra_code);
        let tokens: Vec<(Token, &str)> = collect_tokens(&mut lexer);

        let expected = vec![
            (Token::LessThan, "<"),
            (Token::GreaterThan, ">"),
            (Token::LessThanEqual, "<="),
            (Token::GreaterThanEqual, ">="),
            (Token::DoubleEquals, "=="),
            (Token::Assignment, ":="),
            (Token::Integer, "-42"),
            (Token::Float, "123.45"),
            (Token::Float, "1.2e3"),
            (Token::RadixNumber, "16rA000"),
            (Token::RadixNumber, "2r1010"),
            (Token::String, "'he''llo'"),
            (Token::Symbol, "#'symbol'"),
            (Token::Character, "$x"),
            (Token::Array, "#($a #a 'b' 2 2.0)"),
            (Token::ByteArray, "#[1 2 3]"),
            (Token::Comment, "\"This is a comment\""),
        ];

        assert_token_seq(&tokens, &expected);
    }

    fn collect_tokens<'a>(lexer: &'a mut Lexer<'a>) -> Vec<(Token, &'a str)> {
        let mut tokens = Vec::new();
        while let Some((token, _)) = lexer.next_token() {
            if !matches!(token, Token::Whitespace) {
                tokens.push((token, lexer.slice()));
            }
        }
        tokens
    }

    fn assert_token_seq(actual: &[(Token, &str)], expected: &[(Token, &str)]) {
        println!("--- Tokens received from lexer ---");
        for (i, (token, text)) in actual.iter().enumerate() {
            println!("{:02}: {:?}: '{}'", i, token, text);
        }
        println!("--- End of token list ---\n");

        assert_eq!(
            actual.len(),
            expected.len(),
            "Expected {} tokens but got {}",
            expected.len(),
            actual.len()
        );

        for (i, ((actual_tok, actual_text), (expected_tok, expected_text))) in
            actual.iter().zip(expected.iter()).enumerate()
        {
            assert_eq!(
                actual_tok, expected_tok,
                "Token at position {}: expected {:?} but got {:?}",
                i, expected_tok, actual_tok
            );
            assert_eq!(
                actual_text, expected_text,
                "Text at position {}: expected '{}' but got '{}'",
                i, expected_text, actual_text
            );
        }
    }
}
