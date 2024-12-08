package parser

import (
	"fmt"
	"strings"

	"github.com/neet-007/glox/pkg/scanner"
)

type VisitStmt interface {
	VisitReturnStmt(stmt Return) any
	VisitFunctionStmt(stmt Function) any
	VisitVarDeclaration(stmt VarDeclaration) any
	VisitWhileStmt(stmt WhileStmt) any
	VisitBlockStmt(stmt Block) any
	VisitIfStmt(stmt IfStmt) any
	VisitExpressionStmt(stmt ExpressionStmt) any
	VisitPrintStmt(stmt PrintStmt) any
}

type Stmt interface {
	Accept(visitor VisitStmt) any
}

type Return struct {
	Keyword scanner.Token
	Value   Expr
}

func NewReturn(keyword scanner.Token, value Expr) Return {
	return Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (r Return) String() string {
	return fmt.Sprintf("return keyword:\n %v value %v:\n\n", r.Keyword, r.Value)
}

func (r Return) Accept(visitor VisitStmt) any {
	return visitor.VisitReturnStmt(r)
}

type Function struct {
	Name       scanner.Token
	Parameters []scanner.Token
	Body       []Stmt
}

func NewFunction(name scanner.Token, parameters []scanner.Token, body []Stmt) Function {
	return Function{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}
}

func (f Function) String() string {
	return fmt.Sprintf("name:%v\n paramters:%v\n, body:%v\n\n", f.Name, f.Parameters, f.Body)
}

func (f Function) Accept(visitor VisitStmt) any {
	return visitor.VisitFunctionStmt(f)
}

type VarDeclaration struct {
	Initizlier Expr
	Name       scanner.Token
}

func (v VarDeclaration) String() string {
	return fmt.Sprintf("init:%v name:%v\n", v.Initizlier, v.Name)
}

func NewVarDeclaration(name scanner.Token, initizlier Expr) VarDeclaration {
	return VarDeclaration{
		Name:       name,
		Initizlier: initizlier,
	}
}

func (v VarDeclaration) Accept(visitor VisitStmt) any {
	return visitor.VisitVarDeclaration(v)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("body:%v conditno:%v\n", w.Body, w.Condition)
}

func NewWhileStmt(condition Expr, block Stmt) WhileStmt {
	return WhileStmt{
		Condition: condition,
		Body:      block,
	}
}

func (w WhileStmt) Accept(visitor VisitStmt) any {
	return visitor.VisitWhileStmt(w)
}

type Block struct {
	Statements []Stmt
}

func (b Block) String() string {
	var builder strings.Builder
	builder.WriteString("block:\n")
	for _, stmt := range b.Statements {
		builder.WriteString(fmt.Sprintf("%v\n", stmt))
	}

	return fmt.Sprintf("%s\n", builder.String())
}

func NewBlock(statements []Stmt) Block {
	return Block{
		Statements: statements,
	}
}

func (b Block) Accept(visitor VisitStmt) any {
	return visitor.VisitBlockStmt(b)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i IfStmt) String() string {
	return fmt.Sprintf("condtion:%v then:%v else:%v\n", i.Condition, i.ThenBranch, i.ElseBranch)
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) IfStmt {
	return IfStmt{condition, thenBranch, elseBranch}
}

func (i IfStmt) Accept(visitor VisitStmt) any {
	return visitor.VisitIfStmt(i)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e ExpressionStmt) String() string {
	return fmt.Sprintf("expr:%v\n", e.Expression)
}

func NewExpressionStmt(expr Expr) ExpressionStmt {
	return ExpressionStmt{
		Expression: expr,
	}
}

func (e ExpressionStmt) Accept(visitor VisitStmt) any {
	return visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (p PrintStmt) String() string {
	return fmt.Sprintf("print expr:%v\n", p.Expression)
}

func NewPrintStmt(expr Expr) PrintStmt {
	return PrintStmt{
		Expression: expr,
	}
}

func (p PrintStmt) Accept(visitor VisitStmt) any {
	return visitor.VisitPrintStmt(p)
}
