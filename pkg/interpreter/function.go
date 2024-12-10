package interpreter

import (
	"fmt"

	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/runtime"
)

type LoxFunction struct {
	closure      *runtime.Environment
	Declaration  parser.Function
	isInitilizer bool
}

func NewLoxFunction(stmt parser.Function, closure *runtime.Environment, isIntitlizer bool) LoxFunction {
	return LoxFunction{
		closure:      closure,
		Declaration:  stmt,
		isInitilizer: isIntitlizer,
	}
}

func (l LoxFunction) Call(interpreter *Interpreter, arguemnts []any) (any, error) {
	enviroemnt := runtime.NewEnvironment(l.closure)

	for i := range l.Declaration.Parameters {
		enviroemnt.Define(l.Declaration.Parameters[i].Lexeme, arguemnts[i])
	}

	err := interpreter.executeBlock(l.Declaration.Body, enviroemnt)
	if err != nil {
		if returnVal, ok := err.(*runtime.Return); ok {
			if l.isInitilizer {
				return l.closure.GetAt(0, "this")
			}
			return returnVal.Value, nil
		}

		fmt.Printf("call func 3 err : %v (type: %T)\n", err, err)
		return nil, err
	}

	if l.isInitilizer {
		return l.closure.GetAt(0, "this")
	}
	return nil, nil
}

func (l LoxFunction) Arity() int {
	return len(l.Declaration.Parameters)
}

func (l LoxFunction) Bind(instance Instance) LoxFunction {
	environment := runtime.NewEnvironment(l.closure)
	environment.Define("this", instance)
	return NewLoxFunction(l.Declaration, environment, l.isInitilizer)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.Declaration.Name.Lexeme)
}
