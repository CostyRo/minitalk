use crate::lexer::{Lexer, Token};

pub struct Repl;

impl Repl {
    pub fn new() -> Self {
        Repl
    }

    pub fn start(&self) {
        use std::io::{self, Write};

        loop {
            print!(">>> ");
            io::stdout().flush().unwrap();

            let mut input = String::new();
            if io::stdin().read_line(&mut input).is_err() {
                println!("Failed to read input.");
                continue;
            }

            let input = input.trim();
            if input == "exit" {
                break;
            }

            let mut lexer = Lexer::new(input);
            while let Some((token, _)) = lexer.next_token() {
                if !matches!(token, Token::Whitespace) {
                    println!("{:?}: '{}'", token, lexer.slice());
                }
            }
        }
    }
}
