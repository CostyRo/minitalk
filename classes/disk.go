package classes

import (
	"os"
	"path/filepath"

	"minitalk/types"
	"minitalk/types/core"
)

func NewDiskClass() *core.Object {
	obj := core.NewObject("", "Disk")

	obj.Set("referenceTo", func(args core.Object) interface{} {
		if args.Class != "String" {
			return nil
		}
		rawPath := args.Self.(string)
		cleanPath := filepath.Clean(rawPath)
		return *NewFileClass(cleanPath)
	})

	obj.Set("ls", func(args core.Object) interface{} {
		if args.Class != "String" {
			return nil
		}
		rawPath := args.Self.(string)

		var targetDir string
		if rawPath == "." {
			cwd, err := os.Getwd()
			if err != nil {
				return types.NewArrayObject([]*core.Object{}).Object
			}
			targetDir = cwd
		} else {
			targetDir = filepath.Clean(rawPath)
		}

		entries, err := os.ReadDir(targetDir)
		if err != nil {
			return types.NewArrayObject([]*core.Object{}).Object
		}

		var arr []*core.Object
		for _, e := range entries {
			arr = append(arr, &types.NewStringObject(e.Name()).Object)
		}
		return types.NewArrayObject(arr).Object
	})

	return obj
}
