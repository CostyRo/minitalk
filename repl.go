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

func (r *Repl) processLine(tokens []Token) []core.Object {
	var results []core.Object

	var subTokens []Token
	var stack []core.Object
	var lastType string
	var lastMessage any
	var binaryMessage any
	var lastVar string
	assigment := false
	typeError := false
	messageError := false
	sign := false
	paren := false
	if len(tokens) == 1 && tokens[0].Type == RParen {
		return nil
	}
	for _, tok := range tokens {
		if tok.Type == Whitespace {
			continue
		}

		if paren {
			subTokens = append(subTokens, tok)
			switch tok.Type {
			case RParen:
				if !paren {
					subTokens = subTokens[:len(subTokens)-1]
				}
			case LParen:
				paren = true
				continue
			default:
				continue
			}
		}

		switch tok.Type {
		case LParen:
			paren =  true

		case RParen:
			if paren {
				subResult := r.processLine(subTokens)
				paren = false
				if len(subResult) == 0 {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
					return nil
				}

				obj := subResult[len(subResult)-1]
				if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
					result := fn(obj)
					objResult, ok := result.(core.Object)
					if !ok {
						err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and %s", lastType, obj.Class))
						stack = append(stack, err.Object)
						continue
					}
					stack = append(stack, objResult)
					lastMessage = nil
				} else {
					stack = append(stack, obj)
				}
			}
			
		case Period:
			if len(stack) > 0 {
        		result := stack[len(stack)-1]
        		results = append(results, result)
    		}

			stack = nil
			lastMessage = nil
			lastType = ""
			lastVar = ""
			typeError = false
			messageError = false
			sign = false
			assigment = false

		case Assignment:
			if lastVar == "" {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			assigment = true

		case Colon:
			lastMessage = binaryMessage

		case Symbol, Character, String, Integer, Float, RadixNumber, True, False, Nil:
			var typeName string
			var obj core.Object

			switch tok.Type {
			case Symbol:
				typeName = "Symbol"
				obj = types.NewSymbolObject(tok.Value).Object
			case Character:
				typeName = "Character"
				val := tok.Value[1:]
				obj = types.NewCharacterObject([]rune(val)[0]).Object
			case String:
				typeName = "String"
				val := tok.Value[1:len(tok.Value)-1]
				obj = types.NewStringObject(val).Object
			case Integer:
				typeName = "Integer"
				value, _ := strconv.ParseInt(tok.Value, 10, 64)
				if sign {
					value = -value
					sign = false
				}
				obj = types.NewIntegerObject(value).Object
			case Float:
				typeName = "Float"
				value, _ := strconv.ParseFloat(tok.Value, 64)
				if sign {
					value = -value
					sign = false
				}
				obj = types.NewFloatObject(value).Object
			case RadixNumber:
				typeName = "Integer"
				parts := strings.Split(tok.Value, "r")
				base, _ := strconv.ParseInt(parts[0], 10, 32)
				num, err := strconv.ParseInt(parts[1], int(base), 64)
				if base < 2 || base > 36 {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid base", base)
					return nil
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid number in base", base)
					return nil
				}
				if sign {
					num = -num
					sign = false
				}
				obj = types.NewIntegerObject(num).Object
			case True:
				typeName = "Bool"
				obj = types.NewBoolObject(true).Object
			case False:
				typeName = "Bool"
				obj = types.NewBoolObject(false).Object
			case Nil:
				typeName = "Nil"
				nilObj := core.NewObject(nil, "Nil")
				obj = *nilObj
			}

			if typeError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and %s", lastType, typeName))
				stack = append(stack, err.Object)
				continue
			}
			if messageError {
				err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
				stack = append(stack, err.Object)
				continue
			}

			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(obj)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and %s", lastType, typeName))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				if binaryMessage != nil {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
					return nil
				}
				if assigment {
					r.globalScope[lastVar] = obj
					lastVar = ""
					assigment = false
				} else {
					stack = append(stack, obj)
				}
			}

		case ByteArray:
			valStr := tok.Value[2:len(tok.Value)-1]
			innerTokens := Lex(valStr)

			var elements []byte
			valid := true
			minus := false

			for _, t := range innerTokens {
				var intVal int64
				var err error

				switch t.Type {
				case Whitespace, Plus:
					continue

				case Minus:
					minus = true
					continue

				case Integer:
					intVal, err = strconv.ParseInt(t.Value, 10, 64)

				case RadixNumber:
					parts := strings.Split(t.Value, "r")
					base, _ := strconv.ParseInt(parts[0], 10, 32)
					intVal, err = strconv.ParseInt(parts[1], int(base), 64)
					if base < 2 || base > 36 {
						fmt.Fprintln(os.Stderr, "SyntaxError: invalid base", base)
						return nil
					}
					if err != nil {
						fmt.Fprintln(os.Stderr, "SyntaxError: invalid number in base", base)
						return nil
					}

				case Float:
					var floatVal float64
					floatVal, err = strconv.ParseFloat(t.Value, 64)
					intVal = int64(floatVal)
				
				case Character:
					runes := []rune(t.Value[1:])
					if len(runes) != 1 {
						err = fmt.Errorf("invalid literal: %s", t.Value)
					} else {
						intVal = int64(runes[0])
					}

				case String:
					strVal := t.Value[1:len(t.Value)-1]
					intVal, err = strconv.ParseInt(strVal, 10, 64)
					if err != nil {
						err = fmt.Errorf("invalid literal: %s", t.Value)
					}

				case Symbol:
					symName := t.Value[1:]
					intVal, err = strconv.ParseInt(symName, 10, 64)
					if err != nil {
						err = fmt.Errorf("invalid literal: %s", t.Value)
					}

				case True:
					intVal = 1
				case False:
					intVal = 0
				
				case Identifier:
					obj, inScope := r.globalScope[t.Value]
					if inScope {
						ok := false
						val, ok := obj.Get("toInteger")
						if !ok {
							err = fmt.Errorf("invalid literal: %s", t.Value)
						}
						if intVal, ok = val.(int64); !ok {
							err = fmt.Errorf("invalid literal: %s", t.Value)
						}
					} else {
						valid = false
						stack = append(stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", t.Value)).Object)
					}

				default:
					err = fmt.Errorf("invalid byte array element: %s", t.Value)
				}

				if err != nil || intVal < 0 || intVal > 255 || minus {
					value := t.Value
					if minus {
						value = fmt.Sprintf("-%s", t.Value)
					}
					stack = append(stack, errors.NewValueError(fmt.Sprintf("Invalid byte value: %s", value)).Object)
					valid = false
					break
				}
				elements = append(elements, byte(intVal))
			}

			if !valid {
				continue
			}
			obj := types.NewByteArrayObject(elements).Object
			if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
				result := fn(obj)
				objResult, ok := result.(core.Object)
				if !ok {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and ByteArray", lastType))
					stack = append(stack, err.Object)
					continue
				}
				stack = append(stack, objResult)
				lastMessage = nil
			} else {
				if binaryMessage != nil {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
					return nil
				}
				if assigment {
					r.globalScope[lastVar] = obj
					lastVar = ""
					assigment = false
				} else {
					stack = append(stack, obj)
				}
			}

		case Identifier:
			obj, inScope := r.globalScope[tok.Value]

			if lastMessage != nil{
				if inScope {
					if fn, ok := lastMessage.(func(core.Object) interface{}); ok {
						result := fn(obj)
						objResult, ok := result.(core.Object)
						if !ok {
							err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and %s", lastType, obj.Class))
							stack = append(stack, err.Object)
							continue
						}
						stack = append(stack, objResult)
						lastMessage = nil
					}
					continue
				} else {
					err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s and NameError", lastType))
					stack = append(stack, err.Object)
					continue
				}
			} else if inScope && binaryMessage != nil {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}

			if len(stack) == 0 {
				if inScope {
					stack = append(stack, obj)
				}
				lastVar = tok.Value
			} else {
				last := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				lastType = last.Class
				val, ok := last.Get(tok.Value)
				if !ok {
					messageError = true
					continue
				}
				if fn, ok := val.(func(core.Object) interface{}); ok{
					binaryMessage = fn
				} else if obj, ok := val.(core.Object); ok {
					stack = append(stack, obj)
				} else {
					if constructor := last.GetPropertyType(tok.Value); constructor != nil {
						if obj := constructor(val); obj != nil {
							stack = append(stack, *obj)
						}
					} else {
						stack = append(stack, *types.ObjectConstructor(val))
					}
				}
			}

		case Plus, Minus, Star, Slash, Ampersand, LessThan, GreaterThan, LessThanEqual, GreaterThanEqual, DoubleEquals:
			opMethods := map[TokenType]string{
				Plus:            "plus",
				Minus:           "minus",
				Star:            "mul",
				Slash:           "div",
				Ampersand:       "and",
				LessThan:        "lt",
				GreaterThan:     "gt",
				LessThanEqual:   "le",
				GreaterThanEqual:"ge",
				DoubleEquals:    "eq",
			}

			if _, ok := r.globalScope[lastVar]; lastVar != "" && !ok {
				stack = append(stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", lastVar)).Object)
				lastVar = ""
			}

			if tok.Type == Plus || tok.Type == Minus {
				if tok.Type == Plus {
					if len(stack) == 0 {
						sign = false
						continue
					}
				} else {
					if lastMessage != nil {
						sign = !sign
						continue
					}
					if len(stack) == 0 {
						sign = !sign
						continue
					}
				}
			} else {
				if len(stack) == 0 {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
					return nil
				}
			}

			last := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lastType = last.Class

			val, ok := last.Get(opMethods[tok.Type])
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

	if _, ok := r.globalScope[lastVar]; lastVar != "" && !ok {
		stack = append(stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", lastVar)).Object)
	}

	if len(stack) > 0 {
    	result := stack[len(stack)-1]
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

		tokens := Lex(input)
		outputs := r.processLine(tokens)
		for _, out := range outputs {
    		fmt.Println(out.String())
		}
	}
}
