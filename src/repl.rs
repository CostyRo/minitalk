use std::any::Any;
use std::sync::Arc;

use rustyline::{ Editor, Config, history::MemHistory };

use crate::object::Object;
use crate::float::Float;
use crate::integer::Integer;
use crate::lexer::{ Lexer, Token };

pub struct Repl {
    global_scope: std::collections::HashMap<String, Object>,
}

impl Repl {
    pub fn new() -> Self {
        Self {
            global_scope: std::collections::HashMap::new(),
        }
    }

    fn process_line(&mut self, line: &str) -> Option<Object> {
        let mut stack: Vec<Object> = Vec::new();
        let mut last_message: Option<Box<dyn Any>> = None;
        let mut sign = false;
        let mut lexer = Lexer::new(line);

        while let Some((token, _span)) = lexer.next_token() {
            if matches!(token, Token::Integer) {
                let value = lexer.slice().to_string().parse::<i64>().unwrap();
                let value = if sign { -1 * value } else { value };
                let int = Integer::new(value);
                sign = false;

                if let Some(last_msg) = &last_message {
                    if
                        let Some(func) = last_msg.downcast_ref::<
                            Arc<dyn (Fn(&Box<dyn Any>) -> Object) + Send + Sync>
                        >()
                    {
                        let operand = Box::new(value) as Box<dyn Any>;
                        let result_obj = func(&operand);
                        stack.push(result_obj);
                        last_message = None;
                    } else {
                        panic!("Error!");
                    }
                } else {
                    stack.push(int.obj);
                }
            } else if matches!(token, Token::Float) {
                let value = lexer.slice().to_string().parse::<f64>().unwrap();
                let value = if sign { -1.0 * value } else { value };
                let float = Float::new(value);
                sign = false;

                if let Some(last_msg) = &last_message {
                    if
                        let Some(func) = last_msg.downcast_ref::<
                            Arc<dyn (Fn(&Box<dyn Any>) -> Object) + Send + Sync>
                        >()
                    {
                        let operand = Box::new(value) as Box<dyn Any>;
                        let result_obj = func(&operand);
                        stack.push(result_obj);
                        last_message = None;
                    } else {
                        panic!("Error!");
                    }
                } else {
                    stack.push(float.obj);
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
                            let value = if sign { -1 * value } else { value };
                            let int = Integer::new(value);
                            sign = false;

                            if let Some(last_msg) = &last_message {
                                if
                                    let Some(func) = last_msg.downcast_ref::<
                                        Arc<dyn (Fn(&Box<dyn Any>) -> Object) + Send + Sync>
                                    >()
                                {
                                    let operand = Box::new(value) as Box<dyn Any>;
                                    let result_obj = func(&operand);
                                    stack.push(result_obj);
                                    last_message = None;
                                } else {
                                    panic!("Error!");
                                }
                            } else {
                                stack.push(int.obj);
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
                if let Some(last_obj) = stack.pop() {
                    if
                        let Some(add_func) = last_obj.get::<
                            Arc<dyn (Fn(&Box<dyn Any>) -> Object) + Send + Sync>
                        >("add")
                    {
                        let add_func_cloned = add_func.clone();
                        last_message = Some(Box::new(add_func_cloned) as Box<dyn Any>);
                    } else {
                        println!("No 'add' function found");
                    }
                } else {
                    sign = false;
                }
            } else if matches!(token, Token::Minus) {
                if last_message.is_some(){
                    sign = !sign;
                    continue;
                }
                if let Some(last_obj) = stack.pop(){
                    if
                        let Some(sub_func) = last_obj.get::<
                            Arc<dyn (Fn(&Box<dyn Any>) -> Object) + Send + Sync>
                        >("sub")
                    {
                        let sub_func_cloned = sub_func.clone();
                        last_message = Some(Box::new(sub_func_cloned) as Box<dyn Any>);
                    } else {
                        println!("No 'sub' function found");
                    }
                } else {
                    sign = !sign;
                }
            } else if matches!(token, Token::Star) {
                if let Some(last_obj) = stack.pop() {
                    if
                        let Some(mul_func) = last_obj.get::<
                            Arc<dyn (Fn(&Box<dyn Any>) -> Object) + Send + Sync>
                        >("mul")
                    {
                        let mul_func_cloned = mul_func.clone();
                        last_message = Some(Box::new(mul_func_cloned) as Box<dyn Any>);
                    } else {
                        println!("No 'mul' function found");
                    }
                } else {
                    sign = false;
                }
            }
        }

        stack.pop()
    }

    pub fn start(&mut self) {
        let config = Config::builder().build();
        let mut rl = Editor::<(), MemHistory>::with_history(config, MemHistory::new()).unwrap();

        loop {
            let readline = rl.readline(">>> ");
            match readline {
                Ok(line) => {
                    let input = line.trim();
                    if input == "exit" {
                        break;
                    }
                    let _ = rl.add_history_entry(input);

                    if let Some(last) = self.process_line(input) {
                        println!("{}", last);
                    }
                }
                Err(_) => {
                    break;
                }
            }
        }
    }
}
