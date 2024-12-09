package utils

import (
	"fmt"
	"strings"

	"github.com/neet-007/glox/pkg/parser"
)

type AstPrinter struct{}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}

func (a *AstPrinter) VisitAssignExpr(expr parser.Assign) (any, error) {
	return a.parenthesize(fmt.Sprintf("assign %v", expr.Lexem.Lexeme), expr.Expr), nil
}

func (a *AstPrinter) VisitBinaryExpr(expr parser.Binary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a *AstPrinter) VisitCallExpr(expr parser.Call) (any, error) {
	return a.parenthesize("call", append([]parser.Expr{expr.Callee}, expr.Arguments...)...), nil
}

/*
func (a *AstPrinter) VisitGetExpr(expr Get) (any, error) {
	return a.parenthesize(fmt.Sprintf("get %v", expr.Name.Lexeme), expr.Object)
}
*/

func (a *AstPrinter) VisitGroupingExpr(expr parser.Grouping) (any, error) {
	return a.parenthesize("group", expr.Expr), nil
}

func (a *AstPrinter) VisitLiteralExpr(expr parser.Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a *AstPrinter) VisitLogicalExpr(expr parser.Logical) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

/*
	func (a *AstPrinter) VisitSetExpr(expr Set) (any, error) {
		return a.parenthesize(fmt.Sprintf("set %v", expr.Name.Lexeme), expr.Object, expr.Value)
	}

	func (a *AstPrinter) VisitSuperExpr(expr Super) (any, error) {
		return fmt.Sprintf("(super %v)", expr.Method.Lexeme)
	}

	func (a *AstPrinter) VisitThisExpr(expr This) (any, error) {
		return "this"
	}
*/
func (a *AstPrinter) VisitUnaryExpr(expr parser.Unary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (a *AstPrinter) VisitPrintStmt(stmt parser.PrintStmt) (any, error) {
	return fmt.Sprintf("(print %s)", a.parenthesize("value", stmt.Expression)), nil
}

func (a *AstPrinter) VisitExpressionStmt(stmt parser.ExpressionStmt) (any, error) {
	return a.parenthesize("expression", stmt.Expression), nil
}

func (a *AstPrinter) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	if stmt.ElseBranch != nil {
		return fmt.Sprintf("(if %s %s %s)", a.parenthesize("condition", stmt.Condition), a.print(stmt.ThenBranch), a.print(stmt.ElseBranch)), nil
	}
	return fmt.Sprintf("(if %s %s)", a.parenthesize("condition", stmt.Condition), a.print(stmt.ThenBranch)), nil
}

func (a *AstPrinter) VisitBlockStmt(stmt parser.Block) (any, error) {
	var stmts []string
	for _, statement := range stmt.Statements {
		stmts = append(stmts, a.print(statement))
	}
	return fmt.Sprintf("(block %s)", strings.Join(stmts, " ")), nil
}

func (a *AstPrinter) VisitWhileStmt(stmt parser.WhileStmt) (any, error) {
	return fmt.Sprintf("(while %s %s)", a.parenthesize("condition", stmt.Condition), a.print(stmt.Body)), nil
}

func (a *AstPrinter) VisitVarDeclaration(stmt parser.VarDeclaration) (any, error) {
	if stmt.Initizlier != nil {
		return fmt.Sprintf("(var %s %s)", stmt.Name.Lexeme, a.parenthesize("initializer", stmt.Initizlier)), nil
	}
	return fmt.Sprintf("(var %s)", stmt.Name.Lexeme), nil
}

func (a *AstPrinter) VisitVariableExpr(expr parser.Variable) (any, error) {
	return expr.Name.Lexeme, nil
}

func (a *AstPrinter) VisitFunctionStmt(stmt parser.Function) (any, error) {
	var params []string
	for _, param := range stmt.Parameters {
		params = append(params, param.Lexeme)
	}
	bodyStatms := ""
	for _, bodyStmt := range stmt.Body {
		bodyStatms += a.print(bodyStmt)
	}
	return fmt.Sprintf("(fun %s (%s) %s)", stmt.Name.Lexeme, strings.Join(params, " "), bodyStatms), nil
}

func (a *AstPrinter) VisitReturnStmt(stmt parser.Return) (any, error) {
	if stmt.Value != nil {
		return fmt.Sprintf("(return %s)", a.parenthesize("value", stmt.Value)), nil
	}
	return "(return)", nil
}

/*
	func (a *AstPrinter) VisitClassStmt(stmt Class) (any, error) {
		superclass := stmt.Superclass
		var methods []string
		for _, method := range stmt.Methods {
			methods = append(methods, a.print(method))
		}
		return fmt.Sprintf("(class %s superclass [%s] %s)", stmt.Name.Lexeme, a.VisitVariableExpr(superclass), strings.Join(methods, " "))
	}

*/

func (a *AstPrinter) Print(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		fmt.Printf("%v\n", a.print(stmt))
	}
}

func (a *AstPrinter) print(stmt parser.Stmt) string {
	if stmt == nil {
		return "nil"
	}

	val, err := stmt.Accept(a)
	if err != nil {
		panic(err)
	}
	if val == nil {
		return "nil"
	}

	valStrting, ok := val.(string)
	if !ok {
		panic("not ok")
	}

	return valStrting
}

func (a *AstPrinter) parenthesize(name string, exprs ...parser.Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		val, err := expr.Accept(a)
		if err != nil {
			panic(err)
		}
		if val == nil {
			builder.WriteString("nil")
			continue
		}
		valStr, ok := val.(string)
		if !ok {
			builder.WriteString("COULD_NOT_GET_STRING")
			continue
		}
		builder.WriteString(valStr)
	}
	builder.WriteString(")")

	return builder.String()
}
