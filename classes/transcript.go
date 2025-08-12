package classes

import (
	"fmt"
    "strconv"

	"minitalk/types"
	"minitalk/types/core"
)

func NewTranscriptClass() *core.Object {
	obj := core.NewObject("", "Transcript")

	obj.Set("show", func(args core.Object) interface{} {
		if args.Class != "String" {
			return nil
		}
        raw := args.Self.(string)
        unescaped, err := strconv.Unquote(`"` + raw + `"`)
        if err != nil {
            unescaped = raw
        }
		n, _ := fmt.Print(unescaped)
        nBytes := types.NewIntegerObject(int64(n)).Object
        nBytes.Set("!printable",false)
		return nBytes
	})

	return obj
}
