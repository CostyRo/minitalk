package core

import (
	"fmt"
	"strings"
)

type Object struct {
	Self  interface{}
	properties map[string]interface{}
	Class      string
}

func NewObject(Self interface{}, Class string) *Object {
	obj := &Object{
		Self:       Self,
		properties: make(map[string]interface{}),
		Class:      Class,
	}
	obj.Set("isNil", Self == nil)
	return obj
}

func (o *Object) Set(key string, value interface{}) {
	o.properties[key] = value
}

func (o *Object) Get(key string) (interface{}, bool) {
	v, ok := o.properties[key]
	return v, ok
}

func (o *Object) PropertiesLen() int {
	return len(o.properties)
}

func (o *Object) PropertyNames() []string {
	names := make([]string, 0, len(o.properties))
	for k := range o.properties {
		names = append(names, k)
	}
	return names
}

func (o *Object) String() string {
	switch o.Class {
	case "Integer":
		if v, ok := o.Self.(int64); ok {
			return fmt.Sprintf("%d", v)
		}
	case "Float":
		if v, ok := o.Self.(float64); ok {
			return fmt.Sprintf("%.10f", v)
		}
	case "Bool":
		if v, ok := o.Self.(bool); ok {
			return fmt.Sprintf("%t", v)
		}
	case "Symbol":
		if v, ok := o.Self.(string); ok {
			trimmed := strings.Trim(v, "\"'")
			if strings.ContainsAny(trimmed, " \t\n") {
				return fmt.Sprintf("%s'", trimmed)
			}
			return strings.ReplaceAll(trimmed, "'", "")
		}
	case "Character":
		if v, ok := o.Self.(rune); ok {
			return fmt.Sprintf("$%s", string(v))
		}
	case "String":
		if v, ok := o.Self.(string); ok {
			return fmt.Sprintf("'%s'", v)
		}
	default:
		if strings.HasSuffix(o.Class, "Error") {
			if msg, ok := o.Self.(string); ok {
				return fmt.Sprintf("%s: %s", o.Class, msg)
			}
		}
		if o.Self != nil {
			return fmt.Sprintf("<%s at %p>", o.Class, o)
		}
		return "nil"
	}
	return "Error!"
}
