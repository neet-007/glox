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
}

type clockNativeFunction struct{}

func (c clockNativeFunction) Arity() int {
	return 0
}

func (c clockNativeFunction) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixNano()) / 1e9
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
	}
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		err := i.execute(stmt)

		if err != nil {
			//!TODO error
		}
	}
}

func (i *Interpreter) execute(stmt parser.Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

func (i *Interpreter) evaluate(expr parser.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitReturnStmt(stmt parser.Return) (any, error) {
	//fmt.Println("visit return")
	var val any = nil
	var err error
	if stmt.Value != nil {
		val, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	//fmt.Printf("visit return val %v\n", val)
	return nil, runtime.NewReturn(val)
}

func (i *Interpreter) VisitCallExpr(expr parser.Call) (any, error) {
	//fmt.Println("visit call")
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}

	for _, arg := range expr.Arguments {
		argVal, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, argVal)
	}

	callable, ok := callee.(Callable)
	if !ok {
		//!TODO error
		return nil, runtime.NewRuntimeError("not callable")
	}

	if len(arguments) != callable.Arity() {
		//!TODO error
		return nil, runtime.NewRuntimeError(fmt.Sprintf("expect %d parameters got %d arguments", len(arguments), callable.Arity()))
	}

	callVal := callable.Call(i, arguments)

	//fmt.Printf("return val %v\n", callVal)
	return callVal, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt parser.Function) (any, error) {
	//fmt.Println("visit function stmt")
	function := NewLoxFunction(stmt)
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
		return nil, err
	}

	conditionTruthy := i.isTruthy(condition)

	for conditionTruthy {
		stmt.Body.Accept(i)
		condition, err = i.evaluate(stmt.Condition)

		if err != nil {
			return nil, err
		}
		conditionTruthy = i.isTruthy(condition)
	}
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	//fmt.Println("visit if stmt")
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
	//fmt.Println("visit block")
	err := i.executeBlock(stmt.Statements, runtime.NewEnvironment(i.environment))

	if err != nil {
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
		return nil, err
	}

	//fmt.Printf("visit assinge expr %s %v\n", expr.Lexem, val)
	err = i.environment.Assign(expr.Lexem, val)
	if err != nil {
		//!TODO error
		return nil, err
	}

	return val, nil
}

func (i *Interpreter) VisitVariableExpr(expr parser.Variable) (any, error) {
	//fmt.Println("visit var expr")
	val, err := i.environment.Get(expr.Name)
	if err != nil {
		//!TODO error
		return nil, err
	}

	//fmt.Printf("visit var expr %s %v\n", expr.Name.Lexeme, val)
	return val, nil
}

func (i *Interpreter) VisitBinaryExpr(expr parser.Binary) (any, error) {
	//fmt.Println("visit binary expr")
	leftVal, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	rightVal, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("visit binary expr %v %v\n", leftVal, rightVal)
	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil, err
			}

			return left - right, nil

		}
	case scanner.STAR:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil, err
			}

			return left * right, nil
		}
	case scanner.SLASH:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
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

			//!TODO error
			return nil, err
		}
	default:
		{
			//!TODO error
			return nil, nil
		}
	}
}

func (i *Interpreter) VisitLogicalExpr(expr parser.Logical) (any, error) {
	//fmt.Println("visit logicla expr")
	leftVal, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	rightVal, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("visit logical expr %v %v\n", leftVal, rightVal)
	switch expr.Operator.TokenType {
	case scanner.GREATER:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil, err
			}

			return left > right, nil
		}
	case scanner.GREATER_EQUAL:
		{

			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil, err
			}

			return left >= right, nil
		}
	case scanner.LESS:
		{

			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil, err
			}

			return left < right, nil
		}
	case scanner.LESS_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
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

			//!TODO error
			return nil, err
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

			//!TODO error
			return nil, err
		}
	default:
		{
			//!TODO error
			return nil, nil
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
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			rigthNum, err := i.checkNumberOperand(expr.Operator, rigthVal)
			if err != nil {
				//!TODO error
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
			//!TODO error
			return nil, nil
		}
	}
}

func (i *Interpreter) isTruthy(value any) bool {
	if strVal, ok := value.(string); ok {
		return strVal != ""
	}
	if numVal, ok := value.(float64); ok {
		return numVal != 0
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}

	//!TODO error
	return false
}

func (i *Interpreter) checkNumberOperand(operator scanner.Token, operand any) (float64, error) {
	if val, ok := operand.(float64); ok {
		return val, nil
	}

	//!TODO error
	return 0, nil
}

func (i *Interpreter) checkNumberOperands(operator scanner.Token, operandLeft any, operandRight any) (float64, float64, error) {
	if valleft, okLeft := operandLeft.(float64); okLeft {
		if valRight, okRight := operandRight.(float64); okRight {
			return valleft, valRight, nil
		}
	}

	//!TODO error
	return 0, 0, nil
}
