package types

import (
	"minitalk/interfaces"
	"minitalk/tokens"
	"minitalk/types/core"
	"minitalk/types/errors"
)

type CodeBlockObject struct {
	core.Object
}

func NewCodeBlockObject(arguments []string, loc [][][]string,r interfaces.ReplInterface) *CodeBlockObject {
	obj := core.NewObject(nil, "CodeBlock")

	argsObjs := make([]core.Object, len(arguments))
	for i, arg := range arguments {
		argsObjs[i] = NewStringObject(arg).Object
	}
	argsArrayObj := NewArrayObject(argsObjs)
	obj.Set("arguments", argsArrayObj.Object)
	obj.Set("no_arguments", int64(len(arguments)), ObjectConstructor)

	locList := make([]core.Object, len(loc))
	for i, innerTokens := range loc {
		innerList := make([]core.Object, len(innerTokens))
		for j, pair := range innerTokens {
			pairObj := NewArrayObject([]core.Object{
				NewStringObject(pair[0]).Object,
				NewStringObject(pair[1]).Object,
			})
			innerList[j] = pairObj.Object
		}
		locList[i] = NewArrayObject(innerList).Object
	}
	locArrayObj := NewArrayObject(locList)
	obj.Set("loc", locArrayObj.Object)

	obj.Set("value", func(args ...core.Object) interface{} {
		locVal, ok := obj.Get("loc")
		if !ok {
			return core.Object{}
		}
		locObj, ok := locVal.(core.Object)
		if !ok {
			return core.Object{}
		}
		locArr := locObj.Self.([]core.Object)

		argsVal, ok := obj.Get("arguments")
		if !ok {
			return core.Object{}
		}
		argsObj, ok := argsVal.(core.Object)
		if !ok {
			return core.Object{}
		}
		argsSlice, ok := argsObj.Self.([]core.Object)
		if !ok {
			return core.Object{}
		}

		argName := " "

		if len(argsSlice) >= 1 {
			argName, _ = argsSlice[0].Self.(string)
			r.SetVar(":"+argName, args[0])
		}

		var locData [][][]string
		for _, inner := range locArr {
			innerArr := inner.Self.([]core.Object)
			var innerData [][]string
			for _, pairObj := range innerArr {
				pairArr := pairObj.Self.([]core.Object)
				val := pairArr[1].Self.(string)
				if val == argName {
					val = ":" + argName
				}
				innerData = append(innerData, []string{
					pairArr[0].Self.(string),
					val,
				})
			}
			locData = append(locData, innerData)
		}

		if len(argsSlice) > 1 {
			newArgsSlice := argsSlice[1:]

			newCodeBlock := NewCodeBlockObject(
				func() []string {
					res := make([]string, len(newArgsSlice))
					for i, v := range newArgsSlice {
						res[i] = v.Self.(string)
					}
					return res
				}(),
				locData,
				r,
			)

			return newCodeBlock.Object
		}

		var lastResult core.Object
		for _, line := range locData {
			toks := make([]tokens.Token, len(line))
			for i, pair := range line {
				toks[i] = tokens.Token{Type: StringToTokenType(pair[0]), Value: pair[1]}
			}
			results := r.ProcessLine(toks)
			if len(results) > 0 {
				lastResult = results[len(results)-1]
			}
		}

		if len(argsSlice) == 1 {
			for _, name := range r.GetNames() {
				if len(name) > 0 && name[0] == ':' {
					r.DeleteVar(name)
				}
			}
		}

		if lastResult.Self != nil {
			return lastResult
		}

		return core.Object{}
	})
	obj.Set("toInteger", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toFloat", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toBool", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toSymbol", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toCharacter", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toString", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toByteArray", errors.NewTypeError("Invalid conversion to CodeBlock").Object)
	obj.Set("toArray", errors.NewTypeError("Invalid conversion to CodeBlock").Object)

	return &CodeBlockObject{*obj}
}
