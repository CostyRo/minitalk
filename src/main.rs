mod repl;
mod lexer;
mod object;
mod integer;
mod float;

use repl::Repl;
use std::env;

fn compile_file(_filename: &str) {
    println!("compile_file() not yet implemented");
}

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() > 1 {
        let filename = &args[1];
        compile_file(filename);
    } else {
        let mut repl = Repl::new();
        repl.start();
    }
}
