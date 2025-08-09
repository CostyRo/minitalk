package core

import (
	"fmt"
	"strings"
)

type Object struct {
	Self  interface{}
	properties map[string]interface{}
	propertyTypes  map[string]func(interface{}) *Object
	Class      string
}

func NewObject(Self interface{}, Class string) *Object {
	obj := &Object{
		Self:          Self,
		properties:    make(map[string]interface{}),
		propertyTypes: make(map[string]func(interface{}) *Object),
		Class:         Class,
	}
	obj.Set("isNil", Self == nil)
	return obj
}

func (o *Object) Set(key string, value interface{}, constructor ...func(interface{}) *Object) {
	o.properties[key] = value
	if len(constructor) > 0 && constructor[0] != nil {
		o.propertyTypes[key] = constructor[0]
	}
}

func (o *Object) Get(key string) (interface{}, bool) {
	v, ok := o.properties[key]
	return v, ok
}

func (o *Object) GetPropertyType(key string) func(interface{}) *Object {
	if constructor, ok := o.propertyTypes[key]; ok {
		return constructor
	}
	return nil
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
	case "ByteArray":
		if arr, ok := o.Self.([]byte); ok {
			elems := make([]string, len(arr))
			for i, b := range arr {
				elems[i] = fmt.Sprintf("%d", b)
			}
			return "#[" + strings.Join(elems, " ") + "]"
		}
	case "Array":
		switch arr := o.Self.(type) {
		case []Object:
			elems := make([]string, len(arr))
			for i := range arr {
				elems[i] = (&arr[i]).String()
			}
			return "#(" + strings.Join(elems, " ") + ")"
		case []*Object:
			elems := make([]string, len(arr))
			for i, obj := range arr {
				if obj == nil {
					elems[i] = "nil"
				} else {
					elems[i] = obj.String()
				}
			}
			return "#(" + strings.Join(elems, " ") + ")"
		case []interface{}:
			elems := make([]string, len(arr))
			for i, v := range arr {
				switch el := v.(type) {
				case Object:
					elems[i] = (&el).String()
				case *Object:
					if el == nil {
						elems[i] = "nil"
					} else {
						elems[i] = el.String()
					}
				default:
					if v == nil {
						elems[i] = "nil"
					} else {
						elems[i] = fmt.Sprintf("%v", v)
					}
				}
			}
			return "#(" + strings.Join(elems, " ") + ")"
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
