mod lexer;
mod repl;

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
        let repl = Repl::new();
        repl.start();
    }
}
