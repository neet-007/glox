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
	if interpreter.Debug {
		fmt.Printf("function call\n")
	}
	enviroemnt := runtime.NewEnvironment(l.closure)

	for i := range l.Declaration.Parameters {
		enviroemnt.Define(l.Declaration.Parameters[i].Lexeme, arguemnts[i])
	}

	err := interpreter.executeBlock(l.Declaration.Body, enviroemnt)
	if err != nil {
		if interpreter.Debug {
			fmt.Printf("function call err %v %T\n", err, err)
		}
		if returnVal, ok := err.(*runtime.Return); ok {
			if interpreter.Debug {
				fmt.Printf("function call err is for return\n")
			}
			if l.isInitilizer {
				if interpreter.Debug {
					fmt.Printf("function call err is for return for initilizer\n")
				}
				val, err := l.closure.GetAt(0, "this")
				if err != nil {
					if interpreter.Debug {
						fmt.Printf("function call err is for return for initilizer get this err %v %T\n", err, err)
					}
					return nil, err
				}
				if interpreter.Debug {
					fmt.Printf("function call err is for return for initilizer finshed value %v\n", val)
				}
				return val, nil
			}
			if interpreter.Debug {
				fmt.Printf("function call err is for return finished value\n", returnVal.Value)
			}
			return returnVal.Value, nil
		}

		if interpreter.Debug {
			fmt.Printf("function call err is for return err %v %T\n", err, err)
		}
		return nil, err
	}

	if l.isInitilizer {
		if interpreter.Debug {
			fmt.Printf("function call not return stamtent for initilizer\n")
		}
		val, err := l.closure.GetAt(0, "this")
		if err != nil {
			if interpreter.Debug {
				fmt.Printf("function call not return stamtent for initilizer err %v %T\n", err, err)
			}
			return nil, err
		}
		if interpreter.Debug {
			fmt.Printf("function call not return stamtent for initilizer finished value\n", val)
		}
		return val, nil
	}
	if interpreter.Debug {
		fmt.Printf("function call not return stamtent finished\n")
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
