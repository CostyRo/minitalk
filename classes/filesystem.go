package classes

import (
	"minitalk/types/core"
)

func NewFileSystemClass() *core.Object {
	obj := core.NewObject("", "FileSystem")

	obj.Set("disk", *NewDiskClass())

	return obj
}