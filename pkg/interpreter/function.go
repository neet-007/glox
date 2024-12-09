package interpreter

import (
	"fmt"

	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/runtime"
)

type LoxFunction struct {
	Declaration parser.Function
}

func NewLoxFunction(stmt parser.Function) LoxFunction {
	return LoxFunction{
		Declaration: stmt,
	}
}

func (l LoxFunction) Call(interpreter *Interpreter, arguemnts []any) (any, error) {
	enviroemnt := runtime.NewEnvironment(interpreter.environment)

	for i := range l.Declaration.Parameters {
		enviroemnt.Define(l.Declaration.Parameters[i].Lexeme, arguemnts[i])
	}

	err := interpreter.executeBlock(l.Declaration.Body, enviroemnt)
	if err != nil {
		if returnVal, ok := err.(*runtime.Return); ok {
			return returnVal.Value, nil
		}

		//!TODO error
		return nil, err
	}
	return nil, nil
}

func (l LoxFunction) Arity() int {
	return len(l.Declaration.Parameters)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.Declaration.Name.Lexeme)
}
