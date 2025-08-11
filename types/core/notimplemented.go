package core

type NotImplementedObject struct{}

var NotImplemented = NotImplementedObject{}

func (NotImplementedObject) String() string {
	return "NotImplemented"
}
