use std::any::{Any, TypeId};
use std::sync::Arc;

use crate::object::Object;
use crate::float::Float;
use crate::integer::Integer;
use crate::lexer::{Lexer, Token};

pub struct Repl {
    global_scope: std::collections::HashMap<String, Object>,
    stack: Vec<Object>,
    last_message: Option<Box<dyn Any>>,
}

impl Repl {
    pub fn new() -> Self {
        Self {
            global_scope: std::collections::HashMap::new(),
            stack: Vec::new(),
            last_message: None,
        }
    }

    pub fn start(&mut self) {
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
                if matches!(token, Token::Integer) {
                    let value = lexer.slice().to_string().parse::<i64>().unwrap();
                    let int = Integer::new(value);
    
                    if let Some(last_msg) = &self.last_message {
                        if let Some(func) = last_msg.downcast_ref::<Arc<dyn Fn(&Box<dyn Any>) -> Object + Send + Sync>>() {
                            let operand = Box::new(value) as Box<dyn Any>;
                            let result_obj = func(&operand);
    
                            self.stack.push(result_obj);
                            self.last_message = None;
                        } else {
                            panic!("Error!");
                        }
                    } else {
                        self.stack.push(int.obj);
                    }
                } else if matches!(token, Token::Float) {
                    let value = lexer.slice().to_string().parse::<f64>().unwrap();
                    let float = Float::new(value);
    
                    if let Some(last_msg) = &self.last_message {
                        if let Some(func) = last_msg.downcast_ref::<Arc<dyn Fn(&Box<dyn Any>) -> Object + Send + Sync>>() {
                            let operand = Box::new(value) as Box<dyn Any>;
                            let result_obj = func(&operand);

                            self.stack.push(result_obj);
                            self.last_message = None;
                        } else {
                            panic!("Error!");
                        }
                    } else {
                        self.stack.push(float.obj);
                    }
                } else if matches!(token, Token::RadixNumber) {
                    let slice = lexer.slice();
                    if let Some((base_str, digits)) = slice.split_once('r') {
                        if let Ok(base) = u32::from_str_radix(base_str, 10) {
                            if base < 2 || base > 36 {
                                println!("Base {} is out of range. It must be between 2 and 36.", base);
                                continue;
                            }
                            if let Ok(value) = i64::from_str_radix(digits, base) {
                                let int = Integer::new(value);
                
                                if let Some(last_msg) = &self.last_message {
                                    if let Some(func) = last_msg.downcast_ref::<Arc<dyn Fn(&Box<dyn Any>) -> Object + Send + Sync>>() {
                                        let operand = Box::new(value) as Box<dyn Any>;
                                        let result_obj = func(&operand);
                                        self.stack.push(result_obj);
                                        self.last_message = None;
                                    } else {
                                        panic!("Error!");
                                    }
                                } else {
                                    self.stack.push(int.obj);
                                }
                            } else {
                                println!("Invalid number '{}' in base {}", digits, base);
                            }
                        } else {
                            println!("Invalid base: '{}'", base_str);
                        }
                    } else {
                        println!("Invalid radix number format: '{}'", slice);
                    }
                } else if matches!(token, Token::Plus) {
                    if let Some(mut last_obj) = self.stack.pop() {
                        if let Some(add_func) = last_obj.get::<Arc<dyn Fn(&Box<dyn Any>) -> Object + Send + Sync>>("add") {
                            let add_func_cloned = add_func.clone();
                            self.last_message = Some(Box::new(add_func_cloned) as Box<dyn Any>);
                        } else {
                            panic!("No 'add' function found");
                        }
                    }
                }
            }
            if let Some(last) = self.stack.pop() {
                println!("{}", last);
            }
        }
    }
}
