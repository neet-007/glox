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
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(stmt parser.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr parser.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitReturnStmt(stmt parser.Return) any {
	var val any = nil

	if stmt.Value != nil {
		val = i.evaluate(stmt.Value)
	}

	return runtime.NewReturn(val)
}

func (i *Interpreter) VisitCallExpr(expr parser.Call) any {
	callee := i.evaluate(expr.Callee)

	arguments := []any{}

	for _, arg := range expr.Arguments {
		argVal := i.evaluate(arg)
		arguments = append(arguments, argVal)
	}

	callable, ok := callee.(Callable)
	if !ok {
		//!TODO error
		return nil
	}

	if len(arguments) != callable.Arity() {
		//!TODO error
		return nil
	}
	return callable.Call(i, arguments)
}

func (i *Interpreter) VisitFunctionStmt(stmt parser.Function) any {
	function := NewLoxFunction(stmt)
	i.environment.Define(stmt.Name.Lexeme, function)

	return nil
}

func (i *Interpreter) VisitVarDeclaration(stmt parser.VarDeclaration) any {
	var initizlier any
	if stmt.Initizlier != nil {
		initizlier = i.evaluate(stmt.Initizlier)
	}

	i.environment.Define(stmt.Name.Lexeme, initizlier)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt parser.WhileStmt) any {
	condition := i.evaluate(stmt.Condition)
	conditionTruthy := i.isTruthy(condition)

	for conditionTruthy {
		stmt.Body.Accept(i)
		condition = i.evaluate(stmt.Condition)
		conditionTruthy = i.isTruthy(condition)
	}
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt parser.IfStmt) any {
	condition := i.evaluate(stmt.Condition)
	conditionTruthy := i.isTruthy(condition)

	if conditionTruthy {
		stmt.ThenBranch.Accept(i)
	} else {
		stmt.ElseBranch.Accept(i)
	}

	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt parser.Block) any {
	i.executeBlock(stmt.Statements, runtime.NewEnvironment(i.environment))

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt parser.ExpressionStmt) any {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt parser.PrintStmt) any {
	val := i.evaluate(stmt.Expression)
	fmt.Printf("%v\n", val)
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr parser.Assign) any {
	val := i.evaluate(expr.Expr)

	err := i.environment.Assign(expr.Lexem, val)
	if err != nil {
		//!TODO error
		return nil
	}

	return val
}

func (i *Interpreter) VisitVariableExpr(expr parser.Variable) any {
	val, err := i.environment.Get(expr.Name)
	if err != nil {
		//!TODO error
		return nil
	}

	return val
}

func (i *Interpreter) VisitBinaryExpr(expr parser.Binary) any {
	leftVal := i.evaluate(expr.Left)
	rightVal := i.evaluate(expr.Right)
	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left - right

		}
	case scanner.STAR:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left * right
		}
	case scanner.SLASH:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left / right

		}
	case scanner.PLUS:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err == nil {
				return left + right
			}

			if strLeft, ok := leftVal.(string); ok {
				if strRight, ok := rightVal.(string); ok {
					return strLeft + strRight
				}
			}

			//!TODO error
			return nil
		}
	default:
		{
			//!TODO error
			return nil
		}
	}
}

func (i *Interpreter) VisitLogicalExpr(expr parser.Logical) any {
	leftVal := i.evaluate(expr.Left)
	rightVal := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case scanner.GREATER:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left > right
		}
	case scanner.GREATER_EQUAL:
		{

			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left >= right
		}
	case scanner.LESS:
		{

			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left < right
		}
	case scanner.LESS_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err != nil {
				//!TODO error
				return nil
			}

			return left <= right
		}
	case scanner.EQUAL_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err == nil {
				if expr.Operator.TokenType == scanner.EQUAL_EQUAL {
					return left == right
				} else {
					return left != right
				}
			}

			if strLeft, ok := leftVal.(string); ok {
				if strRight, ok := rightVal.(string); ok {
					return strLeft == strRight
				}
			}

			//!TODO error
			return nil
		}
	case scanner.BANG_EQUAL:
		{
			left, right, err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
			if err == nil {
				if expr.Operator.TokenType == scanner.EQUAL_EQUAL {
					return left == right
				} else {
					return left != right
				}
			}

			if strLeft, ok := leftVal.(string); ok {
				if strRight, ok := rightVal.(string); ok {
					return strLeft != strRight
				}
			}

			//!TODO error
			return nil
		}
	default:
		{
			//!TODO error
			return nil
		}
	}
}

func (i *Interpreter) VisitGroupingExpr(expr parser.Grouping) any {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr parser.Literal) any {
	return expr.Value
}

func (i *Interpreter) executeBlock(stmts []parser.Stmt, enviroment *runtime.Environment) error {
	prev := i.environment
	i.environment = enviroment
	for _, stmt_ := range stmts {
		stmt_.Accept(i)
	}

	i.environment = prev
	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr parser.Unary) any {
	rigthVal := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case scanner.MINUS:
		{
			rigthNum, err := i.checkNumberOperand(expr.Operator, rigthVal)
			if err != nil {
				//!TODO error
				return nil
			}
			return -rigthNum

		}
	case scanner.BANG:
		{
			rigthTruthy := i.isTruthy(rigthVal)
			return !rigthTruthy
		}
	default:
		{
			//!TODO error
			return nil
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
