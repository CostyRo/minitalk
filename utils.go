package main

import (
    "fmt"
    "runtime"
)

func Log(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("%s: at %s:%d\n", msg, file, line)
	}
}