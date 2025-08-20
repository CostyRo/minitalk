package core

import "fmt"

type KeywordMessage struct {
	Obj       *Object
	Message   string
	Parameter *Object
	Optionals map[string]*Object
}

func NewKeywordMessage(obj *Object, msg string) *KeywordMessage {
	return &KeywordMessage{
		Obj:       obj,
		Message:   msg,
		Optionals: make(map[string]*Object),
	}
}

func (km *KeywordMessage) SetOptional(key string, obj *Object) {
	if km.Optionals == nil {
		km.Optionals = make(map[string]*Object)
	}
	km.Optionals[key] = obj
}

func (km *KeywordMessage) IsInitialized() bool {
	return km != nil && km.Obj != nil && km.Message != ""
}

func (km *KeywordMessage) ApplyToObject() {
	for key, obj := range km.Optionals {
		km.Obj.SetOptional(km.Message, key, obj)
	}
}

func (km *KeywordMessage) Reset() {
	if km == nil {
		return
	}
	km.Obj = nil
	km.Message = ""
	km.Parameter = nil
	if km.Optionals != nil {
		for k := range km.Optionals {
			km.Optionals[k] = nil
		}
	}
	km.Optionals = make(map[string]*Object)
}

func (km *KeywordMessage) String() string {
	if km == nil {
		return "<nil KeywordMessage>"
	}

	s := "KeywordMessage:\n"
	if km.Obj != nil {
		s += "  Obj: " + km.Obj.String() + "\n"
	} else {
		s += "  Obj: <nil>\n"
	}

	s += "  Message: " + km.Message + "\n"

	if km.Parameter != nil {
		s += "  Parameter: " + km.Parameter.String() + "\n"
	} else {
		s += "  Parameter: <nil>\n"
	}

	if len(km.Optionals) > 0 {
		s += "  Optionals:\n"
		for key, obj := range km.Optionals {
			s += fmt.Sprintf("    %s -> %s\n", key, obj.String())
		}
	} else {
		s += "  Optionals: <empty>\n"
	}

	return s
}
