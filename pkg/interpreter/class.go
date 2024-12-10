package interpreter

type Class struct {
	methods    map[string]LoxFunction
	Name       string
	SuperClass *Class
}

func NewLoxClass(name string, methods map[string]LoxFunction, class *Class) Class {
	return Class{
		Name:       name,
		methods:    methods,
		SuperClass: class,
	}
}

func (c Class) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewInstance(c)

	initilzier, ok := c.FindMethod("init")
	if ok {
		initilzier.Bind(instance).Call(interpreter, arguments)
	}

	return instance, nil
}

func (c Class) Arity() int {
	initilzier, ok := c.FindMethod("init")
	if ok {
		return initilzier.Arity()
	}
	return 0
}

func (c Class) FindMethod(name string) (LoxFunction, bool) {
	method, ok := c.methods[name]
	if !ok && c.SuperClass != nil {
		return c.SuperClass.FindMethod(name)
	}
	return method, ok
}

func (c Class) String() string {
	return c.Name
}
