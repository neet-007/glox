package interpreter

import (
	"fmt"
	"time"

	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/runtime"
	"github.com/neet-007/glox/pkg/scanner"
)

type Interpreter struct {
	globals     *runtime.Environment
	environment *runtime.Environment
	locals      map[parser.Expr]int
	Debug       bool
}

type clockNativeFunction struct{}

func (c clockNativeFunction) Arity() int {
	return 0
}

func (c clockNativeFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

func (c clockNativeFunction) String() string {
	return "<fn native>"
}

func NewInterpreter(debug bool) *Interpreter {
	globals := runtime.NewEnvironment(nil)
	clock := clockNativeFunction{}
	var clockCallabe Callable = clock

	globals.Define("clock", clockCallabe)
	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      map[parser.Expr]int{},
		Debug:       debug,
	}
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) *runtime.RuntimeError {
	for _, stmt := range stmts {
		err := i.execute(stmt)

		if err != nil {
			if runtimeErr, ok := err.(*runtime.RuntimeError); ok {
				return runtimeErr
			}

			panic(fmt.Sprintf("Excpect runtime error got %v", err))
		}
	}

	return nil
}

func (i *Interpreter) ResolveExpr(expr parser.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) execute(stmt parser.Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

func (i *Interpreter) evaluate(expr parser.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitClassStmt(stmt parser.Class) (any, error) {
	if i.Debug {
		fmt.Printf("interpreter visit class name:%s\n", stmt.Name.Lexeme)
	}
	var superClass Class
	var zeroVariabe parser.Variable
	if stmt.SuperClass != zeroVariabe {
		if i.Debug {
			fmt.Printf("interpreter visit class name:%s has superclass\n", stmt.Name.Lexeme)
		}
		superClassVal, err := i.evaluate(stmt.SuperClass)
		if err != nil {
			if i.Debug {
				fmt.Printf("interpreter visit class name:%s superclass err\n", stmt.Name.Lexeme)
			}
			return nil, err
		}

		superClassClass, ok := superClassVal.(Class)
		if !ok {
			if i.Debug {
				fmt.Printf("interpreter visit class name:%s superclass not class\n", stmt.Name.Lexeme)
			}
			return nil, runtime.NewRuntimeError(stmt.Name, "Superclass must be a class")
		}

		superClass = superClassClass
	}
	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.SuperClass != zeroVariabe {
		if i.Debug {
			fmt.Printf("interpreter visit class name:%s make superclass env and define super\n", stmt.Name.Lexeme)
		}
		i.environment = runtime.NewEnvironment(i.environment)
		i.environment.Define("super", superClass)
	}

	methods := map[string]LoxFunction{}

	for _, method := range stmt.Methods {
		methodFunction := NewLoxFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = methodFunction
	}

	class := NewLoxClass(stmt.Name.Lexeme, methods, &superClass)

	if stmt.SuperClass != zeroVariabe {
		i.environment = i.environment.Enclosing
	}

	i.environment.Assign(stmt.Name, class)

	if i.Debug {
		fmt.Printf("interpreter visit class name:%s finished\n", stmt.Name.Lexeme)
	}
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt parser.Return) (any, error) {
	if i.Debug {
		fmt.Printf("interpreter visit return \n")
	}
	var val any = nil
	var err error
	if stmt.Value != nil {
		if i.Debug {
			fmt.Printf("interpreter visit return has value\n")
		}
		val, err = i.evaluate(stmt.Value)
		if err != nil {
			if i.Debug {
				fmt.Printf("interpreter visit return value error %v %T\n", err, err)
			}
			return nil, err
		}
	}

	if i.Debug {
		fmt.Printf("interpreter visit return finish value:%v\n", val)
	}
	return nil, runtime.NewReturn(val)
}

func (i *Interpreter) VisitThisExpr(expr parser.This) (any, error) {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSetExpr(expr parser.Set) (any, error) {
	if i.Debug {
		fmt.Printf("interpreter visit set name:%v\n", expr.Name.Lexeme)
	}
	object, err := i.evaluate(expr.Object)
	if err != nil {
		if i.Debug {
			fmt.Printf("interpreter visit set name:%v object err %v %T\n", expr.Name.Lexeme, err, err)
		}
		return nil, err
	}

	objectInstance, ok := object.(Instance)
	if !ok {
		if i.Debug {
			fmt.Printf("interpreter visit set name:%v not instance\n", expr.Name.Lexeme)
		}
		return nil, runtime.NewRuntimeError(expr.Name, "Only instances have properties")
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		if i.Debug {
			fmt.Printf("interpreter visit set name:%v value err %v %T\n", expr.Name.Lexeme, err, err)
		}
		return nil, err
	}

	objectInstance.Set(expr.Name, value)

	if i.Debug {
		fmt.Printf("interpreter visit set name:%v finished\n", expr.Name.Lexeme)
	}
	return value, nil
}

func (i *Interpreter) VisitGetExpr(expr parser.Get) (any, error) {
	if i.Debug {
		fmt.Printf("interpreter visit get name:%v\n", expr.Name.Lexeme)
	}
	object, err := i.evaluate(expr.Object)
	if err != nil {
		if i.Debug {
			fmt.Printf("interpreter visit get name:%v object err %v %T\n", expr.Name.Lexeme, err, err)
		}
		return nil, err
	}

	if objectInstance, ok := object.(Instance); ok {
		if i.Debug {
			fmt.Printf("interpreter visit get name:%v finished\n", expr.Name.Lexeme)
		}
		return objectInstance.Get(expr.Name)
	}

	if i.Debug {
		fmt.Printf("interpreter visit get name:%v not instance\n", expr.Name.Lexeme)
	}
	return nil, runtime.NewRuntimeError(expr.Name, "Only instances have properties")
}

func (i *Interpreter) VisitCallExpr(expr parser.Call) (any, error) {
	if i.Debug {
		fmt.Printf("interpreter visit call\n")
	}
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		if i.Debug {
			fmt.Printf("interpreter visit call calle error %v %T\n", err, err)
		}
		return nil, err
	}

	arguments := []any{}

	for _, arg := range expr.Arguments {
		argVal, err := i.evaluate(arg)
		if err != nil {
			if i.Debug {
				fmt.Printf("interpreter visit call arg vall error %v %T\n", err, err)
			}
			return nil, err
		}

		arguments = append(arguments, argVal)
	}

	callable, ok := callee.(Callable)
	if !ok {
		if i.Debug {
			fmt.Printf("interpreter visit call not callalbe\n")
		}
		return nil, runtime.NewRuntimeError(expr.Paren, "not callable")
	}

	if len(arguments) != callable.Arity() {
		if i.Debug {
			fmt.Printf("interpreter visit call err args %d vs arity %d\n", len(arguments), callable.Arity())
		}
		return nil, runtime.NewRuntimeError(expr.Paren, fmt.Sprintf("expect %d parameters got %d arguments", callable.Arity(), len(arguments)))
	}

	callVal, tErr := callable.Call(i, arguments)
	if tErr != nil {
		if i.Debug {
			fmt.Printf("interpreter visit call call value err value %v error %v %v\n", callVal, tErr, tErr)
		}
		return nil, tErr
	}

	if i.Debug {
		fmt.Printf("interpreter visit call finished value %v\n", callVal)
	}
	return callVal, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt parser.Function) (any, error) {
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)

	return nil, nil
}

func (i *Interpreter) VisitVarDeclaration(stmt parser.VarDeclaration) (any, error) {
	var initizlier any
	var err error
	if stmt.Initizlier != nil {
		initizlier, err = i.evaluate(stmt.Initizlier)
		if err != nil {
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.Lexeme, initizlier)
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt parser.WhileStmt) (any, error) {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	conditionTruthy := i.isTruthy(condition)

	for conditionTruthy {
		fmt.Println("here")
		_, err = stmt.Body.Accept(i)
		if err != nil {
			return nil, err
		}
		condition, err = i.evaluate(stmt.Condition)

		if err != nil {
			return nil, err
		}
		conditionTruthy = i.isTruthy(condition)
	}
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	condition, err := i.evaluate(stmt.Condition)

	if err != nil {
		return nil, err
	}
	conditionTruthy := i.isTruthy(condition)

	if conditionTruthy {
		_, err = stmt.ThenBranch.Accept(i)
		if err != nil {
			return nil, err
		}
	} else if stmt.ElseBranch != nil {
		_, err = stmt.ElseBranch.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt parser.Block) (any, error) {
	err := i.executeBlock(stmt.Statements, runtime.NewEnvironment(i.environment))

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt parser.ExpressionStmt) (any, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt parser.PrintStmt) (any, error) {
	val, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	if val == nil {
		fmt.Println("nil")
	} else {
		fmt.Printf("%v\n", val)
	}
	return nil, nil
}

func (i *Interpreter) VisitSuperExpr(expr parser.Super) (any, error) {
	dist, ok := i.locals[expr]
	if !ok {
		return nil, runtime.NewRuntimeError(expr.Keyword, "superclass not found")
	}

	class, err := i.environment.GetAt(dist, "super")
	if err != nil {
		return nil, err
	}

	classClass, ok := class.(Class)
	if !ok {
		return nil, runtime.NewRuntimeError(expr.Keyword, "superclass not found")
	}

	instance, err := i.environment.GetAt(dist-1, "this")
	if err != nil {
		return nil, err
	}

	instanceInstance, ok := instance.(Instance)
	if !ok {
		return nil, runtime.NewRuntimeError(expr.Keyword, "instance not found")
	}

	method, ok := classClass.FindMethod(expr.Method.Lexeme)
	if !ok {
		return nil, runtime.NewRuntimeError(expr.Method, "method not found")
	}

	return method.Bind(instanceInstance), nil

}

func (i *Interpreter) VisitAssignExpr(expr parser.Assign) (any, error) {
	val, err := i.evaluate(expr.Expr)

	if err != nil {
		return nil, err
	}

	if dist, ok := i.locals[expr]; ok {
		i.environment.AssignAt(dist, expr.Lexem, val)
	} else {
		tErr := i.globals.Assign(expr.Lexem, val)
		if tErr != nil {
			return nil, tErr
		}
	}

	return val, nil
}

func (i *Interpreter) VisitVariableExpr(expr parser.Variable) (any, error) {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitBinaryExpr(expr parser.Binary) (any, error) {
	leftVal, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	rightVal, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left - right, nil

		}
	case scanner.STAR:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left * right, nil
		}
	case scanner.SLASH:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left / right, nil

		}
	case scanner.PLUS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err == nil {
				return left + right, nil
			}

			if strLeft, ok := leftVal.(string); ok {
				if strRight, ok := rightVal.(string); ok {
					return strLeft + strRight, nil
				}
			}

			return nil, runtime.NewRuntimeError(expr.Operator, "Expect binary operands to be strings")
		}
	case scanner.GREATER:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left > right, nil
		}
	case scanner.GREATER_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left >= right, nil
		}
	case scanner.LESS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left < right, nil
		}
	case scanner.LESS_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				return nil, err
			}

			return left <= right, nil
		}
	case scanner.EQUAL_EQUAL:
		{
			if leftVal == nil && rightVal == nil {
				return true, nil
			}
			if leftVal == nil {
				return false, nil
			}

			return leftVal == rightVal, nil
		}
	case scanner.BANG_EQUAL:
		{
			if leftVal == nil && rightVal == nil {
				return false, nil
			}
			if leftVal == nil {
				return true, nil
			}

			return leftVal != rightVal, nil
		}
	default:
		{
			return nil, runtime.NewRuntimeError(expr.Operator, "Excpect binray operator to be -, +, *, /")
		}
	}
}

func (i *Interpreter) VisitLogicalExpr(expr parser.Logical) (any, error) {
	leftVal, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == scanner.OR {
		if i.isTruthy(leftVal) {
			return leftVal, nil
		}
	} else {
		if !i.isTruthy(leftVal) {
			return leftVal, nil
		}
	}

	rightVal, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	return rightVal, nil
}

func (i *Interpreter) VisitGroupingExpr(expr parser.Grouping) (any, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr parser.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) executeBlock(stmts []parser.Stmt, enviroment *runtime.Environment) error {
	prev := i.environment
	i.environment = enviroment
	for _, stmt_ := range stmts {
		_, err := stmt_.Accept(i)
		if err != nil {
			i.environment = prev
			return err
		}
	}

	i.environment = prev
	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr parser.Unary) (any, error) {
	rigthVal, err := i.evaluate(expr.Right)

	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			rigthNum, err := i.checkNumberOperand(expr.Operator, rigthVal)
			if err != nil {
				return nil, err
			}
			return -rigthNum, nil

		}
	case scanner.BANG:
		{
			rigthTruthy := i.isTruthy(rigthVal)
			return !rigthTruthy, nil
		}
	default:
		{
			return nil, runtime.NewRuntimeError(expr.Operator, "Expect unary operator to be -, !")
		}
	}
}

func (i *Interpreter) lookUpVariable(name scanner.Token, expr parser.Expr) (any, error) {
	if i.Debug {
		fmt.Printf("lookup variable name:%s\n", name.Lexeme)
	}
	if dist, ok := i.locals[expr]; ok {
		if i.Debug {
			fmt.Printf("lookup variable name:%s found dist:%d\n", name.Lexeme, dist)
		}
		val, err := i.environment.GetAt(dist, name.Lexeme)
		if err != nil {
			if i.Debug {
				fmt.Printf("lookup variable name:%s found dist:%d get at error %v %T\n", name.Lexeme, dist, err, err)
			}
			return nil, err
		}

		return val, nil
	} else {
		if i.Debug {
			fmt.Printf("lookup variable name:%s look global\n", name.Lexeme)
		}
		val, err := i.globals.Get(name)
		if err != nil {
			if i.Debug {
				fmt.Printf("lookup variable name:%s global get error %v %T\n", name.Lexeme, dist, err, err)
			}
			return nil, err
		}

		if i.Debug {
			fmt.Printf("lookup variable name:%s finshed val %v\n", name.Lexeme, val)
		}
		return val, nil
	}
}

func (i *Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}

	return true
}

func (i *Interpreter) checkNumberOperand(operator scanner.Token, operand any) (float64, *runtime.RuntimeError) {
	if val, ok := operand.(float64); ok {
		return val, nil
	}

	return 0, runtime.NewRuntimeError(operator, "Expect operands to be numbers")
}

func (i *Interpreter) checkNumberOperands(operator scanner.Token, operandLeft any, operandRight any) (float64, float64, *runtime.RuntimeError) {
	if valleft, okLeft := operandLeft.(float64); okLeft {
		if valRight, okRight := operandRight.(float64); okRight {
			return valleft, valRight, nil
		}
	}

	return 0, 0, runtime.NewRuntimeError(operator, "Expect operands to be numbers")
}
