package runtime

import (
	"github.com/neet-007/glox/pkg/scanner"
)

type Environment struct {
	Enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		values:    map[string]any{},
	}
}

func (e *Environment) Get(name scanner.Token) (any, *RuntimeError) {
	val, ok := e.values[name.Lexeme]
	if !ok {
		if e.Enclosing != nil {
			return e.Enclosing.Get(name)
		}
		return nil, NewRuntimeError("undefiend variable " + name.Lexeme)
	}

	return val, nil
}

func (e *Environment) Assign(name scanner.Token, value any) *RuntimeError {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return NewRuntimeError("undefiend variable " + name.Lexeme)
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}
