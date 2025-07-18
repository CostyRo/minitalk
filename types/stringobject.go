package types

import "minitalk/types/core"

type StringObject struct {
	core.Object
}

func NewStringObject(value string) *StringObject {
	obj := core.NewObject(value, "String")

    obj.Set("add", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewStringObject(value + other.Self.(string)).Object
	})
	obj.Set("lt", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value < other.Self.(string)).Object
	})
	obj.Set("gt", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value > other.Self.(string)).Object
	})
	obj.Set("le", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value <= other.Self.(string)).Object
	})
	obj.Set("ge", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value >= other.Self.(string)).Object
	})
	obj.Set("eq", func(other core.Object) interface{} {
		if other.Class != "String" {
			return nil
		}
		return NewBoolObject(value == other.Self.(string)).Object
	})

	return &StringObject{*obj}
}
