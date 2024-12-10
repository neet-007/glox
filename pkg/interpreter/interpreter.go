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

func NewInterpreter() *Interpreter {
	globals := runtime.NewEnvironment(nil)
	clock := clockNativeFunction{}
	var clockCallabe Callable = clock

	globals.Define("clock", clockCallabe)
	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      map[parser.Expr]int{},
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
	i.environment.Define(stmt.Name.Lexeme, nil)

	methods := map[string]LoxFunction{}

	for _, method := range stmt.Methods {
		methodFunction := NewLoxFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = methodFunction
	}

	class := NewLoxClass(stmt.Name.Lexeme, methods)
	i.environment.Assign(stmt.Name, class)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt parser.Return) (any, error) {
	//fmt.Println("visit return")
	var val any = nil
	var err error
	if stmt.Value != nil {
		val, err = i.evaluate(stmt.Value)
		if err != nil {
			fmt.Printf("returm err: %v %T\n", err, err)
			return nil, err
		}
	}

	//fmt.Printf("visit return val %v\n", val)
	return nil, runtime.NewReturn(val)
}

func (i *Interpreter) VisitThisExpr(expr parser.This) (any, error) {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSetExpr(expr parser.Set) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	objectInstance, ok := object.(Instance)
	if !ok {
		return nil, runtime.NewRuntimeError("Only instances have properties")
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	objectInstance.Set(expr.Name, value)

	return value, nil
}

func (i *Interpreter) VisitGetExpr(expr parser.Get) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		fmt.Printf("get err %v %T\n", err, err)
		return nil, err
	}

	if objectInstance, ok := object.(Instance); ok {
		return objectInstance.Get(expr.Name)
	}

	return nil, runtime.NewRuntimeError("Only instances have properties")
}

func (i *Interpreter) VisitCallExpr(expr parser.Call) (any, error) {
	//fmt.Println("visit call")
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		fmt.Printf("call 1 err: %v %T\n", err, err)
		return nil, err
	}

	arguments := []any{}

	for _, arg := range expr.Arguments {
		argVal, err := i.evaluate(arg)
		if err != nil {
			fmt.Printf("call 2 err: %v %T\n", err, err)
			return nil, err
		}

		arguments = append(arguments, argVal)
	}

	callable, ok := callee.(Callable)
	if !ok {
		return nil, runtime.NewRuntimeError("not callable")
	}

	if len(arguments) != callable.Arity() {
		return nil, runtime.NewRuntimeError(fmt.Sprintf("expect %d parameters got %d arguments", len(arguments), callable.Arity()))
	}

	callVal, tErr := callable.Call(i, arguments)
	if tErr != nil {
		fmt.Printf("call 3 err: %v %T\n", err, err)
		return nil, tErr
	}

	//fmt.Printf("return val %v\n", callVal)
	return callVal, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt parser.Function) (any, error) {
	//fmt.Println("visit function stmt")
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)

	return nil, nil
}

func (i *Interpreter) VisitVarDeclaration(stmt parser.VarDeclaration) (any, error) {
	//fmt.Println("visit var dec")
	var initizlier any
	var err error
	if stmt.Initizlier != nil {
		initizlier, err = i.evaluate(stmt.Initizlier)
		if err != nil {
			fmt.Printf("var dec err: %v %T\n", err, err)
			return nil, err
		}
	}

	//fmt.Printf("visit var dec %s %v\n", stmt.Name.Lexeme, initizlier)
	i.environment.Define(stmt.Name.Lexeme, initizlier)
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt parser.WhileStmt) (any, error) {
	//fmt.Println("visit while")
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		fmt.Printf("while 1  err: %v %T\n", err, err)
		return nil, err
	}

	conditionTruthy, tErr := i.isTruthy(condition)
	if tErr != nil {
		fmt.Printf("while 2  err: %v %T\n", err, err)
		return nil, tErr
	}

	for conditionTruthy {
		stmt.Body.Accept(i)
		condition, err = i.evaluate(stmt.Condition)

		if err != nil {
			fmt.Printf("while 3  err: %v %T\n", err, err)
			return nil, err
		}
		conditionTruthy, tErr = i.isTruthy(condition)
		if tErr != nil {
			fmt.Printf("while 4  err: %v %T\n", err, err)
			return nil, tErr
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	//fmt.Println("visit if stmt")
	condition, err := i.evaluate(stmt.Condition)

	if err != nil {
		fmt.Printf("if 1  err: %v %T\n", err, err)
		return nil, err
	}
	conditionTruthy, tErr := i.isTruthy(condition)
	if tErr != nil {
		fmt.Printf("if 2  err: %v %T\n", err, err)
		return nil, tErr
	}

	if conditionTruthy {
		_, err = stmt.ThenBranch.Accept(i)
		if err != nil {
			fmt.Printf("if 3  err: %v %T\n", err, err)
			return nil, err
		}
	} else if stmt.ElseBranch != nil {
		_, err = stmt.ElseBranch.Accept(i)
		if err != nil {
			fmt.Printf("if 4  err: %v %T\n", err, err)
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt parser.Block) (any, error) {
	//fmt.Println("visit block")
	err := i.executeBlock(stmt.Statements, runtime.NewEnvironment(i.environment))

	if err != nil {
		fmt.Printf("block  1  err: %v %T\n", err, err)
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt parser.ExpressionStmt) (any, error) {
	//fmt.Println("visit expr stmt")
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt parser.PrintStmt) (any, error) {
	//fmt.Println("visit print")
	val, err := i.evaluate(stmt.Expression)
	if err != nil {
		fmt.Printf("print  1  err: %v %T\n", err, err)
		return nil, err
	}
	//fmt.Printf("visit print val %v\n", val)
	fmt.Printf("%v\n", val)
	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr parser.Assign) (any, error) {
	//fmt.Println("visit asggine expr")
	val, err := i.evaluate(expr.Expr)

	if err != nil {
		fmt.Printf("assgin 1  err: %v %T\n", err, err)
		return nil, err
	}

	if dist, ok := i.locals[expr]; ok {
		i.environment.AssignAt(dist, expr.Lexem, val)
	} else {
		//fmt.Printf("visit assinge expr %s %v\n", expr.Lexem, val)
		tErr := i.globals.Assign(expr.Lexem, val)
		if tErr != nil {
			fmt.Printf("assgin 2  err: %v %T\n", err, err)
			return nil, tErr
		}
	}

	return val, nil
}

func (i *Interpreter) VisitVariableExpr(expr parser.Variable) (any, error) {
	//fmt.Println("visit var expr")
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitBinaryExpr(expr parser.Binary) (any, error) {
	//fmt.Println("visit binary expr")
	leftVal, err := i.evaluate(expr.Left)
	if err != nil {
		fmt.Printf("bi 1  err: %v %T\n", err, err)
		return nil, err
	}

	rightVal, err := i.evaluate(expr.Right)
	if err != nil {
		fmt.Printf("bi 2  err: %v %T\n", err, err)
		return nil, err
	}

	//fmt.Printf("visit binary expr %v %v\n", leftVal, rightVal)
	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("bi 3  err: %v %T\n", err, err)
				return nil, err
			}

			return left - right, nil

		}
	case scanner.STAR:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("bi 4  err: %v %T\n", err, err)
				return nil, err
			}

			return left * right, nil
		}
	case scanner.SLASH:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("bi 5  err: %v %T\n", err, err)
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

			return nil, runtime.NewRuntimeError("Expect binary operands to be strings")
		}
	default:
		{
			return nil, runtime.NewRuntimeError("Excpect binray operator to be -, +, *, /")
		}
	}
}

func (i *Interpreter) VisitLogicalExpr(expr parser.Logical) (any, error) {
	//fmt.Println("visit logicla expr")
	leftVal, err := i.evaluate(expr.Left)
	if err != nil {
		fmt.Printf("log 1  err: %v %T\n", err, err)
		return nil, err
	}

	rightVal, err := i.evaluate(expr.Right)
	if err != nil {
		fmt.Printf("log 2  err: %v %T\n", err, err)
		return nil, err
	}

	//fmt.Printf("visit logical expr %v %v\n", leftVal, rightVal)
	switch expr.Operator.TokenType {
	case scanner.GREATER:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("log 3  err: %v %T\n", err, err)
				return nil, err
			}

			return left > right, nil
		}
	case scanner.GREATER_EQUAL:
		{

			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("log 4  err: %v %T\n", err, err)
				return nil, err
			}

			return left >= right, nil
		}
	case scanner.LESS:
		{

			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("log 5  err: %v %T\n", err, err)
				return nil, err
			}

			return left < right, nil
		}
	case scanner.LESS_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				fmt.Printf("log 6  err: %v %T\n", err, err)
				return nil, err
			}

			return left <= right, nil
		}
	case scanner.EQUAL_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err == nil {
				if expr.Operator.TokenType == scanner.EQUAL_EQUAL {
					return left == right, nil
				} else {
					return left != right, nil
				}
			}

			if strLeft, ok := leftVal.(string); ok {
				if strRight, ok := rightVal.(string); ok {
					return strLeft == strRight, nil
				}
			}

			return nil, runtime.NewRuntimeError("Expect binary operands to be strings")
		}
	case scanner.BANG_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err == nil {
				if expr.Operator.TokenType == scanner.EQUAL_EQUAL {
					return left == right, nil
				} else {
					return left != right, nil
				}
			}

			if strLeft, ok := leftVal.(string); ok {
				if strRight, ok := rightVal.(string); ok {
					return strLeft != strRight, nil
				}
			}

			return nil, runtime.NewRuntimeError("Expect binary operands to be strings")
		}
	default:
		{
			return nil, runtime.NewRuntimeError("Excpect logical operator to be >, >=. <, <=, !=, ==")
		}
	}
}

func (i *Interpreter) VisitGroupingExpr(expr parser.Grouping) (any, error) {
	//fmt.Println("visit grpuing")
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr parser.Literal) (any, error) {
	//fmt.Println("visit literal")
	//fmt.Printf("visit literal %v\n", expr.Value)
	return expr.Value, nil
}

func (i *Interpreter) executeBlock(stmts []parser.Stmt, enviroment *runtime.Environment) error {
	//fmt.Println("ecvute block")
	prev := i.environment
	i.environment = enviroment
	for _, stmt_ := range stmts {
		_, err := stmt_.Accept(i)
		if err != nil {
			fmt.Printf("execute block  err: %v %T\n", err, err)
			i.environment = prev
			return err
		}
	}

	i.environment = prev
	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr parser.Unary) (any, error) {
	//fmt.Println("visit unraty ")
	rigthVal, err := i.evaluate(expr.Right)

	if err != nil {
		fmt.Printf("unary 1  err: %v %T\n", err, err)
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			rigthNum, err := i.checkNumberOperand(expr.Operator, rigthVal)
			if err != nil {
				fmt.Printf("unary 2  err: %v %T\n", err, err)
				return nil, err
			}
			return -rigthNum, nil

		}
	case scanner.BANG:
		{
			rigthTruthy, err := i.isTruthy(rigthVal)
			if err != nil {
				fmt.Printf("unary 3  err: %v %T\n", err, err)
				return nil, err
			}
			return !rigthTruthy, nil
		}
	default:
		{
			return nil, runtime.NewRuntimeError("Expect unary operator to be -, !")
		}
	}
}

func (i *Interpreter) lookUpVariable(name scanner.Token, expr parser.Expr) (any, error) {
	if dist, ok := i.locals[expr]; ok {
		val, err := i.environment.GetAt(dist, name.Lexeme)
		if err != nil {
			fmt.Printf("look up  err: %v %T\n", err, err)
			return nil, err
		}

		return val, nil
	} else {
		val, err := i.globals.Get(name)
		if err != nil {
			fmt.Printf("look up2  err: %v %T\n", err, err)
			return nil, err
		}

		return val, nil
	}
}

func (i *Interpreter) isTruthy(value any) (bool, *runtime.RuntimeError) {
	if strVal, ok := value.(string); ok {
		return strVal != "", nil
	}
	if numVal, ok := value.(float64); ok {
		return numVal != 0, nil
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal, nil
	}

	return false, runtime.NewRuntimeError("Unexpected value")
}

func (i *Interpreter) checkNumberOperand(operator scanner.Token, operand any) (float64, *runtime.RuntimeError) {
	if val, ok := operand.(float64); ok {
		return val, nil
	}

	return 0, runtime.NewRuntimeError("Expect operands to be numbers")
}

func (i *Interpreter) checkNumberOperands(operator scanner.Token, operandLeft any, operandRight any) (float64, float64, *runtime.RuntimeError) {
	if valleft, okLeft := operandLeft.(float64); okLeft {
		if valRight, okRight := operandRight.(float64); okRight {
			return valleft, valRight, nil
		}
	}

	return 0, 0, runtime.NewRuntimeError("Expect operands to be numbers")
}
