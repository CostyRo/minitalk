package main

import (
	"fmt"
	"os"
)

func compileFile(filename string) {
	fmt.Println("compileFile() not yet implemented")
}

func main() {
	args := os.Args

	if len(args) > 1 {
		filename := args[1]
		compileFile(filename)
	} else {
		repl := NewRepl()
		repl.Start()
	}
}
