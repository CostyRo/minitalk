package interfaces

import (
	"minitalk/tokens"
	"minitalk/types/core"
)

type ReplInterface interface {
	ProcessLine(tokens []tokens.Token) []core.Object
	GetVar(name string) (core.Object, bool)
	SetVar(name string, val core.Object)
	DeleteVar(name string)
	GetNames() []string
}
