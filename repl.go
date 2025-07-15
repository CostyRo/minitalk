package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peterh/liner"

	"minitalk/types"
)

type Repl struct {
	globalScope map[string]types.Object
	liner       *liner.State
}

func NewRepl() *Repl {
	return &Repl{
		globalScope: make(map[string]types.Object),
		liner:       liner.NewLiner(),
	}
}

func (r *Repl) processLine(input string) *string {
	var stack []types.Object
	var lastMessage any
	sign := false

	tokens := Lex(input)
	for _, tok := range tokens {
		switch tok.Type {
		case Whitespace:
			continue

		case Integer:
			value, _ := strconv.ParseInt(tok.Value, 10, 64)
			if sign {
				value = -value
				sign = false
			}
			intObj := types.NewIntegerObject(value)
			if fn, ok := lastMessage.(func(interface{}) interface{}); ok {
				result := fn(value)
				objResult, ok := result.(types.Object)
				if !ok {
					fmt.Println("function did not return Object")
					return nil
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, intObj.Object)
			}

		case Float:
			value, _ := strconv.ParseFloat(tok.Value, 64)
			if sign {
				value = -value
				sign = false
			}
			floatObj := types.NewFloatObject(value)
			if fn, ok := lastMessage.(func(interface{}) interface{}); ok {
				result := fn(value)
				objResult, ok := result.(types.Object)
				if !ok {
					fmt.Println("function did not return Object")
					return nil
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, floatObj.Object)
			}

		case RadixNumber:
			parts := strings.Split(tok.Value, "r")
			if len(parts) != 2 {
				fmt.Println("invalid radix number format:", tok.Value)
				continue
			}
			base, err1 := strconv.ParseInt(parts[0], 10, 32)
			num, err2 := strconv.ParseInt(parts[1], int(base), 64)
			if err1 != nil || err2 != nil {
				fmt.Println("invalid radix number:", tok.Value)
				continue
			}
			if base < 2 || base > 36 {
				fmt.Printf("base %d out of range\n", base)
				continue
			}
			if sign {
				num = -num
				sign = false
			}
			intObj := types.NewIntegerObject(num)
			if fn, ok := lastMessage.(func(interface{}) interface{}); ok {
				result := fn(num)
				objResult, ok := result.(types.Object)
				if !ok {
					fmt.Println("function did not return Object")
					return nil
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, intObj.Object)
			}

		case Plus:
			if len(stack) == 0 {
				sign = false
				continue
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			val, ok := last.Get("add")
			if !ok {
				fmt.Println("no 'add' function found")
				continue
			}
			fn, ok := val.(func(interface{}) interface{})
			if !ok {
				fmt.Println("'add' is not a valid function")
				continue
			}
			lastMessage = fn

		case Minus:
			if lastMessage != nil {
				sign = !sign
				continue
			}
			if len(stack) == 0 {
				sign = !sign
				continue
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			val, ok := last.Get("sub")
			if !ok {
				fmt.Println("no 'sub' function found")
				continue
			}
			fn, ok := val.(func(interface{}) interface{})
			if !ok {
				fmt.Println("'sub' is not a valid function")
				continue
			}
			lastMessage = fn

		case Star:
			if len(stack) == 0 {
				sign = false
				continue
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			val, ok := last.Get("mul")
			if !ok {
				fmt.Println("no 'mul' function found")
				continue
			}
			fn, ok := val.(func(interface{}) interface{})
			if !ok {
				fmt.Println("'mul' is not a valid function")
				continue
			}
			lastMessage = fn
		}
	}

	if len(stack) > 0 {
		result := stack[len(stack)-1].String()
		return &result
	}
	return nil
}

func (r *Repl) Start() {
	defer r.liner.Close()

	r.liner.SetCtrlCAborts(true)
	r.liner.SetMultiLineMode(false)

	for {
		line, err := r.liner.Prompt(">>> ")
		if err != nil {
			fmt.Println()
			break
		}

		input := strings.TrimSpace(line)
		if input == "exit" {
			break
		}

		r.liner.AppendHistory(input)

		if output := r.processLine(input); output != nil {
			fmt.Println(*output)
		}
	}
}
