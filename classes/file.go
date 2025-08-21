package classes

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"minitalk/types"
	"minitalk/types/core"
)

func NewFileClass(path string) *core.Object {
	p := filepath.Clean(path)
	abs, _ := filepath.Abs(p)
	info, err := os.Stat(p)
	exists := err == nil

	obj := core.NewObject(p, "File")

	obj.Set("basename", types.NewStringObject(filepath.Base(p)).Object)
	obj.Set("baseNameWithoutExtension", types.NewStringObject(strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))).Object)
	obj.Set("extension", types.NewStringObject(filepath.Ext(p)).Object)
	obj.Set("fullName", types.NewStringObject(abs).Object)
	obj.Set("path", types.NewStringObject(p).Object)
	obj.Set("parent", types.NewStringObject(filepath.Dir(p)).Object)
	if exists {
		obj.Set("size", types.NewIntegerObject(info.Size()).Object)
		obj.Set("modificationTime", types.NewStringObject(info.ModTime().Format(time.RFC3339)).Object)
	} else {
		obj.Set("size", types.NewIntegerObject(0).Object)
		obj.Set("modificationTime", types.NewStringObject("").Object)
	}
	obj.Set("exists", types.NewBoolObject(exists).Object)
	obj.Set("isAbsent", types.NewBoolObject(!exists).Object)
	obj.Set("isFile", types.NewBoolObject(exists && info.Mode().IsRegular()).Object)
	obj.Set("isDirectory", types.NewBoolObject(exists && info.IsDir()).Object)
	obj.Set("isSymlink", types.NewBoolObject(exists && isSymlink(p)).Object)
	obj.Set("isEmpty", types.NewBoolObject(isEmptyPath(p)).Object)
	obj.Set("isReadable", types.NewBoolObject(isReadable(p)).Object)
	obj.Set("isWritable", types.NewBoolObject(isWritable(p)).Object)
	obj.Set("isExecutable", types.NewBoolObject(exists && info.Mode()&0111 != 0).Object)
	obj.Set("isHidden", types.NewBoolObject(strings.HasPrefix(filepath.Base(p), ".")).Object)
	obj.Set("contents", func() core.Object {
		data, err := os.ReadFile(p)
		if err != nil {
			return types.NewStringObject("").Object
		}
		return types.NewStringObject(string(data)).Object
	})
	obj.Set("tell", types.NewIntegerObject(0).Object)
	obj.Set("seek", func(arg core.Object) interface{} {
		if arg.Class != "Integer" {
			return nil
		}
		obj.Set("tell", arg)
		return types.NewBoolObject(true).Object
	})
	obj.Set("nextLine", func() core.Object {
		f, err := os.Open(p)
		if err != nil {
			return types.NewStringObject("").Object
		}
		defer f.Close()

		var pos int64
		tell, ok := obj.Get("tell")
		if ok {
			if tellObj, ok := tell.(core.Object); ok {
				pos = tellObj.Self.(int64)
			}
		}

		_, err = f.Seek(pos, io.SeekStart)
		if err != nil {
			return types.NewStringObject("").Object
		}

		reader := bufio.NewReader(f)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return types.NewStringObject("").Object
		}

		pos += int64(len(line))
		obj.Set("tell", types.NewIntegerObject(pos).Object)

		return types.NewStringObject(strings.TrimRight(line, "\r\n")).Object
	})
	obj.Set("write", func(arg core.Object) interface{} {
		if arg.Class != "String" {
			return nil
		}
		content := arg.Self.(string)
		err := os.WriteFile(p, []byte(content), 0644)
		if err != nil {
			return types.NewBoolObject(false).Object
		}
		info, err := os.Stat(p)
		if err == nil {
			obj.Set("size", types.NewIntegerObject(info.Size()).Object)
			obj.Set("modificationTime", types.NewStringObject(info.ModTime().Format(time.RFC3339)).Object)
		}
		return types.NewBoolObject(true).Object
	})

	obj.Set("append", func(arg core.Object) interface{} {
		if arg.Class != "String" {
			return nil
		}
		content := arg.Self.(string)
		f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return types.NewBoolObject(false).Object
		}
		defer f.Close()

		_, err = f.WriteString(content)
		if err != nil {
			return types.NewBoolObject(false).Object
		}

		info, err := os.Stat(p)
		if err == nil {
			obj.Set("size", types.NewIntegerObject(info.Size()).Object)
			obj.Set("modificationTime", types.NewStringObject(info.ModTime().Format(time.RFC3339)).Object)
		}

		return types.NewBoolObject(true).Object
	})

	return obj
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

func isEmptyPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		f, _ := os.Open(path)
		names, _ := f.Readdirnames(1)
		f.Close()
		return len(names) == 0
	}
	return info.Size() == 0
}

func isReadable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func isWritable(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	f.Close()
	return true
}
