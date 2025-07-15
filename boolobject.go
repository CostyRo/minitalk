package main

type BoolObject struct {
	Object
}

func NewBoolObject(value bool) *BoolObject {
	obj := NewObject(value, "bool")
	return &BoolObject{*obj}
}
