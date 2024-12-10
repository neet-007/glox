package interpreter

import (
	"github.com/neet-007/glox/pkg/runtime"
	"github.com/neet-007/glox/pkg/scanner"
)

type Instance struct {
	class  Class
	fields map[string]any
}

func NewInstance(class Class) Instance {
	return Instance{
		class:  class,
		fields: map[string]any{},
	}
}

func (i Instance) Get(name scanner.Token) (any, error) {
	if val, ok := i.fields[name.Lexeme]; ok {
		return val, nil
	}

	method, ok := i.class.FindMethod(name.Lexeme)
	if ok {
		return method.Bind(i), nil
	}

	return nil, runtime.NewRuntimeError("Undefined property '" + name.Lexeme)
}

func (i Instance) Set(name scanner.Token, value any) (any, error) {
	i.fields[name.Lexeme] = value
	return nil, nil
}

func (i Instance) String() string {
	return i.class.Name + " instance"
}
