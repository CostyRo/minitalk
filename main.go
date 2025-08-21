package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"minitalk/tokens"
)

func lineFeeder(lines []string) func(string) (string, error) {
	i := 0
	return func(string) (string, error) {
		if i >= len(lines) {
			return "", io.EOF
		}
		s := lines[i]
		i++
		return s, nil
	}
}

func compileFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	repl := NewRepl()
	handler := NewInputHandler()
	lines := strings.Split(string(data), "\r\n")
	feed := lineFeeder(lines)

	for {
		line, err := feed(">>> ")
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		input, err := handler.Complete(line, feed)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			break
		}
		toks := filterWhitespace(tokens.Lex(input))
		repl.ProcessLine(toks)
	}
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
