package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/peterh/liner"

	"minitalk/tokens"
	"minitalk/types"
	"minitalk/types/core"
	"minitalk/types/errors"
)

type Repl struct {
	globalScope map[string]core.Object
	liner       *liner.State
}

func (r *Repl) GetVar(name string) (core.Object, bool) {
	val, ok := r.globalScope[name]
	return val, ok
}

func (r *Repl) SetVar(name string, val core.Object) {
	if r.globalScope == nil {
		r.globalScope = make(map[string]core.Object)
	}
	r.globalScope[name] = val
}

func (r *Repl) DeleteVar(name string) {
	delete(r.globalScope, name)
}

func (r *Repl) GetNames() []string {
	names := make([]string, 0, len(r.globalScope))
	for name := range r.globalScope {
		names = append(names, name)
	}
	return names
}

func NewRepl() *Repl {
	return &Repl{
		globalScope: make(map[string]core.Object),
		liner:       liner.NewLiner(),
	}
}

func parseByteArray(value string, r *Repl, stack *[]core.Object) ([]byte, bool) {
	valStr := value[2:len(value)-1]
	innerTokens := tokens.Lex(valStr)

	var elements []byte
	valid := true
	minus := false

	for _, t := range innerTokens {
		var intVal int64
		var err error

		switch t.Type {
		case tokens.Whitespace, tokens.Plus:
			continue

		case tokens.Minus:
			minus = true
			continue

		case tokens.Integer:
			intVal, err = strconv.ParseInt(t.Value, 10, 64)

		case tokens.RadixNumber:
			parts := strings.Split(t.Value, "r")
			base, _ := strconv.ParseInt(parts[0], 10, 32)
			intVal, err = strconv.ParseInt(parts[1], int(base), 64)
			if base < 2 || base > 36 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid base", base)
				return nil, false
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid number in base", base)
				return nil, false
			}

		case tokens.Float:
			var floatVal float64
			floatVal, err = strconv.ParseFloat(t.Value, 64)
			intVal = int64(floatVal)

		case tokens.Character:
			runes := []rune(t.Value[1:])
			if len(runes) != 1 {
				err = fmt.Errorf("invalid literal: %s", t.Value)
			} else {
				intVal = int64(runes[0])
			}

		case tokens.String:
			strVal := t.Value[1:len(t.Value)-1]
			intVal, err = strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				err = fmt.Errorf("invalid literal: %s", t.Value)
			}

		case tokens.Symbol:
			symName := t.Value[1:]
			intVal, err = strconv.ParseInt(symName, 10, 64)
			if err != nil {
				err = fmt.Errorf("invalid literal: %s", t.Value)
			}

		case tokens.True:
			intVal = 1
		case tokens.False:
			intVal = 0

		case tokens.Identifier:
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
				*stack = append(*stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", t.Value)).Object)
			}

		default:
			err = fmt.Errorf("invalid byte array element: %s", t.Value)
		}

		if err != nil || intVal < 0 || intVal > 255 || minus {
			valueStr := t.Value
			if minus {
				valueStr = fmt.Sprintf("-%s", t.Value)
			}
			*stack = append(*stack, errors.NewValueError(fmt.Sprintf("Invalid byte value: %s", valueStr)).Object)
			valid = false
			break
		}

		elements = append(elements, byte(intVal))
		minus = false
	}

	return elements, valid
}

func parseArray(value string, r *Repl, stack *[]core.Object) ([]core.Object, bool) {
	valStr := value[2:len(value)-1]
	innerTokens := tokens.Lex(valStr)

	var elements []core.Object
	valid := true
	minus := false

	for _, t := range innerTokens {
		var obj core.Object
		var err error

		switch t.Type {
		case tokens.Whitespace, tokens.Plus:
			continue

		case tokens.Minus:
			minus = true
			continue

		case tokens.Integer:
			intVal, _ := strconv.ParseInt(t.Value, 10, 64)
			if minus {
				intVal = -intVal
			}
			obj = types.NewIntegerObject(intVal).Object

		case tokens.RadixNumber:
			parts := strings.Split(t.Value, "r")
			base, _ := strconv.ParseInt(parts[0], 10, 32)
			var intVal int64
			intVal, err = strconv.ParseInt(parts[1], int(base), 64)
			if base < 2 || base > 36 {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid base", base)
				return nil, false
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid number in base", base)
				return nil, false
			}
			if minus {
				intVal = -intVal
			}
			obj = types.NewIntegerObject(intVal).Object

		case tokens.Float:
			var floatVal float64
			floatVal, _ = strconv.ParseFloat(t.Value, 64)
			if minus {
				floatVal = -floatVal
			}
			obj = types.NewFloatObject(floatVal).Object

		case tokens.Character:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Character")
				return nil, false
			}
			obj = types.NewCharacterObject([]rune(t.Value[1:])[0]).Object

		case tokens.String:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for String")
				return nil, false
			}
			obj = types.NewStringObject(t.Value[1:len(t.Value)-1]).Object

		case tokens.Symbol:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Symbol")
				return nil, false
			}
			obj = types.NewSymbolObject(t.Value[1:]).Object

		case tokens.True:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Bool")
				return nil, false
			}
			obj = types.NewBoolObject(true).Object

		case tokens.False:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Character")
				return nil, false
			}
			obj = types.NewBoolObject(false).Object

		case tokens.Nil:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Nil")
				return nil, false
			}
			obj = *core.NewObject(nil, "Nil")

		case tokens.Identifier:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for variables")
				return nil, false
			}
			var inScope bool
			if obj, inScope = r.globalScope[t.Value]; !inScope {
				valid = false
				*stack = append(*stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", t.Value)).Object)
			}

		case tokens.ByteArray:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for ByteArray")
				return nil, false
			}
			bytes, ok := parseByteArray(t.Value, r, stack)
			if ok {
				obj = types.NewByteArrayObject(bytes).Object
			} else {
				err = fmt.Errorf("invalid bytearray literal: %s", t.Value)
			}

		case tokens.Array:
			if minus {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Array")
				return nil, false
			}
			nested, ok := parseArray(t.Value, r, stack)
			if ok {
				obj = types.NewArrayObject(nested).Object
			} else {
				fmt.Println("A")
				err = fmt.Errorf("invalid array literal: %s", t.Value)
			}

		default:
			err = fmt.Errorf("invalid array element: %s", t.Value)
		}

		if err != nil {
			*stack = append(*stack, errors.NewValueError(err.Error()).Object)
			valid = false
			break
		}

		elements = append(elements, obj)
		minus = false
	}

	return elements, valid
}

func (r *Repl) ProcessLine(toks []tokens.Token) []core.Object {
	var results []core.Object

	var subTokens []tokens.Token
	var stack []core.Object
	var argumentsCodeBlock []string
	var lastType string
	var lastMessage any
	var binaryMessage any
	var lastVar string
	locCodeBlock := [][][]string{{}}
	assigment := false
	typeError := false
	messageError := false
	minus := false
	plus := false
	paren := false
	bracket := 0
	argumentCodeBlock := false
	pipe := false
	nonPipe := false

	if len(toks) == 1 && toks[0].Type == tokens.RParen {
		return nil
	}

	for _, tok := range toks {
		plus = false
		if tok.Type == tokens.Whitespace {
			continue
		}

		if paren {
			subTokens = append(subTokens, tok)
			switch tok.Type {
			case tokens.RParen:
				if !paren {
					subTokens = subTokens[:len(subTokens)-1]
				}
			case tokens.LParen:
				paren = true
				continue
			default:
				continue
			}
		}

		if bracket != 0 {
			switch tok.Type {
			case tokens.LBracket:
				bracket++
				if pipe {
					locCodeBlock[len(locCodeBlock)-1] = append(
						locCodeBlock[len(locCodeBlock)-1],
						[]string{TokenTypeToString(tok.Type), tok.Value})
				}
			case tokens.RBracket:
				bracket--
				if bracket == 0 {
					pipe = false
					obj := types.NewCodeBlockObject(argumentsCodeBlock, locCodeBlock, r)
					stack = append(stack, obj.Object)
				} else if pipe {
					locCodeBlock[len(locCodeBlock)-1] = append(
						locCodeBlock[len(locCodeBlock)-1],
						[]string{TokenTypeToString(tok.Type), tok.Value})
				}
			case tokens.Colon:
				argumentCodeBlock = true
			case tokens.Identifier:
				if pipe {
					locCodeBlock[len(locCodeBlock)-1] = append(
						locCodeBlock[len(locCodeBlock)-1],
						[]string{TokenTypeToString(tok.Type), tok.Value})
				} else if argumentCodeBlock {
					argumentsCodeBlock = append(argumentsCodeBlock, tok.Value)
					argumentCodeBlock = false
				}
			case tokens.Pipe:
				pipe = true
				if nonPipe {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid characters in arguments list of a code block")
					return nil
				}
			case tokens.Period:
				if pipe {
					locCodeBlock = append(locCodeBlock, [][]string{})
				} else {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid characters in arguments list of a code block")
					return nil
				}
			default:
				locCodeBlock[len(locCodeBlock)-1] = append(
					locCodeBlock[len(locCodeBlock)-1],
					[]string{TokenTypeToString(tok.Type), tok.Value})
				nonPipe = true
			}
			continue
		}

		switch tok.Type {
		case tokens.Illegal:
			fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
			return nil

		case tokens.LParen:
			paren = true

		case tokens.RParen:
			if paren {
				subResult := r.ProcessLine(subTokens)
				paren = false
				if len(subResult) == 0 {
					fmt.Fprintln(os.Stderr, "SyntaxError: empty parenthesis")
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

		case tokens.LBracket:
			bracket = 1

		case tokens.Period:
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
			minus = false
			plus = false
			assigment = false

		case tokens.Assignment:
			if lastVar == "" {
				fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
				return nil
			}
			assigment = true

		case tokens.Colon:
			lastMessage = binaryMessage

		case tokens.Symbol, tokens.Character, tokens.String, tokens.Integer, tokens.Float,
			tokens.RadixNumber, tokens.True, tokens.False, tokens.Nil:

			var typeName string
			var obj core.Object

			switch tok.Type {
			case tokens.Symbol:
				if minus {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Symbol")
					return nil
				}
				typeName = "Symbol"
				obj = types.NewSymbolObject(tok.Value).Object

			case tokens.Character:
				if minus {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Character")
					return nil
				}
				typeName = "Character"
				val := tok.Value[1:]
				obj = types.NewCharacterObject([]rune(val)[0]).Object

			case tokens.String:
				if minus {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for String")
					return nil
				}
				typeName = "String"
				val := tok.Value[1 : len(tok.Value)-1]
				obj = types.NewStringObject(val).Object

			case tokens.Integer:
				typeName = "Integer"
				value, _ := strconv.ParseInt(tok.Value, 10, 64)
				if minus {
					value = -value
					minus = false
				}
				obj = types.NewIntegerObject(value).Object

			case tokens.Float:
				typeName = "Float"
				value, _ := strconv.ParseFloat(tok.Value, 64)
				if minus {
					value = -value
					minus = false
				}
				obj = types.NewFloatObject(value).Object

			case tokens.RadixNumber:
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
				if minus {
					num = -num
					minus = false
				}
				obj = types.NewIntegerObject(num).Object

			case tokens.True:
				if minus {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Bool")
					return nil
				}
				typeName = "Bool"
				obj = types.NewBoolObject(true).Object

			case tokens.False:
				if minus {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Bool")
					return nil
				}
				typeName = "Bool"
				obj = types.NewBoolObject(false).Object

			case tokens.Nil:
				if minus {
					fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for Nil")
					return nil
				}
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
			} else if fn, ok := lastMessage.(func(...core.Object) interface{}); ok {
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
					fmt.Println("FOUND!")
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

		case tokens.ByteArray, tokens.Array:
			var typeName string
			var obj core.Object

			if tok.Type == tokens.ByteArray {
				typeName = "ByteArray"
				elements, ok := parseByteArray(tok.Value, r, &stack)
				if !ok {
					continue
				}
				obj = types.NewByteArrayObject(elements).Object
			} else {
				typeName = "Array"
				elements, ok := parseArray(tok.Value, r, &stack)
				if !ok {
					continue
				}
				obj = types.NewArrayObject(elements).Object
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

		case tokens.Identifier:
			obj, inScope := r.globalScope[tok.Value]

			if lastMessage != nil {
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
			} else if len(stack) == 0 {
				if inScope {
					if minus {
						fmt.Fprintln(os.Stderr, "SyntaxError: invalid unary minus for variables")
						return nil
					}
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
				if fn, ok := val.(func(core.Object) interface{}); ok {
					binaryMessage = fn
				} else if fnCodeBlock, ok := val.(func(...core.Object) interface{}); ok {
					if noArgs, ok := last.Get("no_arguments"); ok {
						if valInt, ok := noArgs.(int64); ok {
							if valInt == 0 {
								objResult, ok := fnCodeBlock().(core.Object)
								if !ok{
									Log("COMPILER ERROR!")
									return nil
								}
								stack = append(stack, objResult)
							} else {
								binaryMessage = fnCodeBlock
							}
						}
					} else {
						Log("COMPILER ERROR!")
						return nil
					}
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

		case tokens.Plus, tokens.Minus, tokens.Star, tokens.Slash, tokens.Ampersand,
			tokens.LessThan, tokens.GreaterThan, tokens.LessThanEqual, tokens.GreaterThanEqual, tokens.DoubleEquals:
			opMethods := map[tokens.TokenType]string{
				tokens.Plus:             "plus",
				tokens.Minus:            "minus",
				tokens.Star:             "mul",
				tokens.Slash:            "div",
				tokens.Ampersand:        "and",
				tokens.LessThan:         "lt",
				tokens.GreaterThan:      "gt",
				tokens.LessThanEqual:    "le",
				tokens.GreaterThanEqual: "ge",
				tokens.DoubleEquals:     "eq",
			}

			if _, ok := r.globalScope[lastVar]; lastVar != "" && !ok {
				stack = append(stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", lastVar)).Object)
				lastVar = ""
			}

			if tok.Type == tokens.Plus || tok.Type == tokens.Minus {
				if tok.Type == tokens.Plus {
					plus = true
					if len(stack) == 0 {
						minus = false
						continue
					}
				} else {
					if lastMessage != nil {
						minus = !minus
						continue
					}
					if len(stack) == 0 {
						minus = !minus
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
		if len(stack) == 0 {
			stack = append(stack, errors.NewNameError(fmt.Sprintf("'%s' is not defined", lastVar)).Object)
		} else {
			r.globalScope[lastVar] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
	}

	if messageError {
		err := errors.NewTypeError(fmt.Sprintf("Message doesn't exists for %s", lastType))
		stack = append(stack, err.Object)
	}

	if len(stack) > 0 {
		result := stack[len(stack)-1]
		r.globalScope["_"] = result
		results = append(results, result)
	} else if minus || plus {
		fmt.Fprintln(os.Stderr, "SyntaxError: invalid syntax")
		return nil
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

		toks := tokens.Lex(input)
		outputs := r.ProcessLine(toks)
		for _, out := range outputs {
			fmt.Println(out.String())
		}
	}
}
