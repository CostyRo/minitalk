package classes

import (
	"minitalk/global"
	"minitalk/types"
	"minitalk/types/core"
)

func NewStdinClass() *core.Object {
	obj := core.NewObject("", "Stdin")

	obj.Set("nextLine", func() core.Object {
		line, err := global.Liner.Prompt("")
		if err != nil {
			return types.NewStringObject("").Object
		}
		return types.NewStringObject(line).Object
	})

	return obj
}
