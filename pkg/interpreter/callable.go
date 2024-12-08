package interpreter

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) any
	String() string
}
