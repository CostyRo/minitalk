package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/peterh/liner"

	"minitalk/types"
	"minitalk/types/core"
	"minitalk/types/errors"
)

type Repl struct {
	globalScope map[string]core.Object
	liner       *liner.State
}

func NewRepl() *Repl {
	return &Repl{
		globalScope: make(map[string]core.Object),
		liner:       liner.NewLiner(),
	}
}

func (r *Repl) processLine(input string) []string {
	var results []string

	var stack []core.Object
	var lastType string
	var lastMessage any
	typeError := false
	messageError := false
	sign := false

	tokens := Lex(input)
	for _, tok := range tokens {
		switch tok.Type {
		case Whitespace:
			continue

		case Period:
			if len(stack) > 0 {
        		result := stack[len(stack)-1].String()
        		results = append(results, result)
    		}

			stack = nil
			lastMessage = nil
			lastType = ""
			typeError = false
			messageError = false
			sign = false

		case Integer:
			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Integer", lastType))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}
			value, _ := strconv.ParseInt(tok.Value, 10, 64)
			if sign {
				value = -value
				sign = false
			}
			intObj := types.NewIntegerObject(value)
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(intObj.Object)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Integer", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, intObj.Object)
			}

		case Float:
			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Float", lastType))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}
			value, _ := strconv.ParseFloat(tok.Value, 64)
			if sign {
				value = -value
				sign = false
			}
			floatObj := types.NewFloatObject(value)
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(floatObj.Object)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Float", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, floatObj.Object)
			}

		case RadixNumber:
			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Integer", lastType))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}
			parts := strings.Split(tok.Value, "r")
			base, _ := strconv.ParseInt(parts[0], 10, 32)
			num, err := strconv.ParseInt(parts[1], int(base), 64)
			if base < 2 || base > 36 {
				fmt.Fprintln(os.Stderr, "Syntax error: invalid base", base)
				return nil
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "Syntax error: invalid number in base", base)
				return nil
			}
			if sign {
				num = -num
				sign = false
			}
			intObj := types.NewIntegerObject(num)
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(intObj.Object)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Integer", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, intObj.Object)
			}

		case True:
			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Bool", lastType))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}
			boolObj := types.NewBoolObject(true)
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(boolObj.Object)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Bool", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, boolObj.Object)
			}

		case False:
			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Bool", lastType))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}
			boolObj := types.NewBoolObject(false)
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(boolObj.Object)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Bool", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, boolObj.Object)
			}

		case Nil:
			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Nil", lastType))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}
			nilObj := core.NewObject(nil, "Nil")
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(*nilObj)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and Nil", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				stack = append(stack, *nilObj)
			}

		case Identifier:
			if len(stack) == 0 {
				continue
			} else {
				last := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				lastType = last.Class
				_, ok := last.Get(tok.Value)
				if !ok {
					messageError = true
					continue
				}
			}

		case Plus:
			if len(stack) == 0 {
				sign = false
				continue
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("add")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
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
			lastType = last.Class
			val, ok := last.Get("sub")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case Star:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("mul")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case Slash:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("div")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn
		
		case Ampersand:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("and")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case LessThan:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("lt")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case GreaterThan:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("gt")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case LessThanEqual:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("le")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case GreaterThanEqual:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("ge")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn

		case DoubleEquals:
			if len(stack) == 0 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class
			val, ok := last.Get("eq")
			if !ok {
				typeError = true
				continue
			}
			fn, ok := val.(func(core.Object) interface{})
			if !ok {
				Log("COMPILER ERROR")
				return nil
			}
			lastMessage = fn	
		}
	}

	if len(stack) > 0 {
    	result := stack[len(stack)-1].String()
    	results = append(results, result)
	}
	return results
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

		outputs := r.processLine(input)
		for _, out := range outputs {
    		fmt.Println(out)
		}
	}
}
