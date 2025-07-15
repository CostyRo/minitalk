package types

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type Object struct {
	selfValue  interface{}
	properties map[string]interface{}
	class      string
}

func NewObject(selfValue interface{}, class string) *Object {
	return &Object{
		selfValue:  selfValue,
		properties: make(map[string]interface{}),
		class:      class,
	}
}

func (o *Object) SetSelfValue(value interface{}, class string) {
	o.selfValue = value
	o.class = class
}

func GetSelfValue[T any](o *Object) (T, bool) {
	val, ok := o.selfValue.(T)
	return val, ok
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
	switch o.class {
	case "integer":
		if v, ok := o.selfValue.(int64); ok {
			return fmt.Sprintf("%d", v)
		}
	case "float":
		if v, ok := o.selfValue.(float64); ok {
			return fmt.Sprintf("%.10f", v)
		}
	case "bool":
		if v, ok := o.selfValue.(bool); ok {
			return fmt.Sprintf("%t", v)
		}
	default:
		if strings.HasSuffix(o.class, "error") {
			if msg, ok := o.selfValue.(string); ok {
				return fmt.Sprintf("%s: %s", o.class, msg)
			}
		}
		if o.selfValue != nil {
			ptr := reflect.ValueOf(o.selfValue).Pointer()
			return fmt.Sprintf("<%s at %p>", o.class, unsafe.Pointer(ptr))
		}
		return "nil"
	}
	return "Error!"
}
