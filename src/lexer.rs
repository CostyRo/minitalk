use logos::{Logos, Span};

#[derive(Logos, Debug, PartialEq)]
pub enum Token {
    // Keywords
    #[token("self")]
    Self_,
    #[token("super")]
    Super,
    #[token("nil")]
    Nil,
    #[token("true")]
    True,
    #[token("false")]
    False,

    // Symbols
    #[token("(")]
    LParen,
    #[token(")")]
    RParen,
    #[token("[")]
    LBracket,
    #[token("]")]
    RBracket,
    #[token(".")]
    Period,
    #[token(";")]
    Semicolon,
    #[token(":")]
    Colon,
    #[token("|")]
    Pipe,
    #[token("^")]
    Caret,

    // Operators
    #[token("+")]
    Plus,
    #[token("-")]
    Minus,
    #[token("*")]
    Star,
    #[token("/")]
    Slash,
    #[token("&")]
    Ampersand,
    #[token("<")]
    LessThan,
    #[token(">")]
    GreaterThan,
    #[token("<=")]
    LessThanEqual,
    #[token(">=")]
    GreaterThanEqual,
    #[token("==")]
    DoubleEquals,
    #[token(":=")]
    Assignment,

    #[regex("[a-zA-Z][a-zA-Z0-9_]*")]
    Identifier,

    // Example: 123, -42
    #[regex(r"[+-]?[0-9]+")]
    Integer,

    // Example: 123.45, -1.2e3, 4.5E+6
    #[regex(r"[+-]?(?:[0-9]+\.[0-9]+(?:[eE][+-]?[0-9]+)?|[0-9]+(?:[eE][+-]?[0-9]+))")]
    Float,

    // Example: 2r1010, 16rA000
    #[regex(r"[+-]?[0-9]+r[0-9A-Fa-f]+")]
    RadixNumber,

    // Example: 'hello', 'it''s fine'
    #[regex(r"'([^']|'{2})*'")]
    String,

    #[regex(r"#'([^']|'{2})*'|#[a-zA-Z_][a-zA-Z0-9_]*")]
    Symbol,

    #[regex("\\$.")]
    Character,

    #[regex("#\\([^)]*\\)")]
    Array,

    #[regex("#\\[[^\\]]*\\]")]
    ByteArray,

    #[regex("\"[^\"]*\"")]
    Comment,

    #[regex("[ \t\r\n]+")]
    Whitespace,

    Error,
}

pub struct Lexer<'source> {
    logos_lexer: logos::Lexer<'source, Token>,
}

impl<'source> Lexer<'source> {
    pub fn new(input: &'source str) -> Self {
        Self {
            logos_lexer: Token::lexer(input),
        }
    }

    pub fn next_token(&mut self) -> Option<(Token, Span)> {
        match self.logos_lexer.next() {
            Some(Ok(token)) => Some((token, self.logos_lexer.span())),
            Some(Err(_)) => Some((Token::Error, self.logos_lexer.span())),
            None => None,
        }
    }

    pub fn slice(&self) -> &'source str {
        self.logos_lexer.slice()
    }
}
