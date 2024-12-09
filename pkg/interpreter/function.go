package interpreter

import (
	"fmt"

	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/runtime"
)

type LoxFunction struct {
	closure     *runtime.Environment
	Declaration parser.Function
}

func NewLoxFunction(stmt parser.Function, closure *runtime.Environment) LoxFunction {
	return LoxFunction{
		closure:     closure,
		Declaration: stmt,
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
			return returnVal.Value, nil
		}

		fmt.Printf("call func 3 err : %v (type: %T)\n", err, err)
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
