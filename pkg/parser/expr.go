package parser

import (
	"github.com/neet-007/glox/pkg/scanner"
)

type VisitExpr interface {
	VisitVariableExpr(expr Variable) any
	VisitAssignExpr(expr Assign) any
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitLogicalExpr(expr Logical) any
	VisitUnaryExpr(expr Unary) any
}

type Expr interface {
	Accept(visitor VisitExpr) any
}

type Variable struct {
	Name scanner.Token
}

func NewVariable(name scanner.Token) Variable {
	return Variable{
		Name: name,
	}
}

func (v Variable) Accept(visitor VisitExpr) any {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Lexem scanner.Token
	Expr  Expr
}

func NewAssign(lexem scanner.Token, expr Expr) Assign {
	return Assign{
		Lexem: lexem,
		Expr:  expr,
	}
}

func (a Assign) Accept(visitor VisitExpr) any {
	return visitor.VisitAssignExpr(a)
}

type Binary struct {
	Left     Expr
	Right    Expr
	Operator scanner.Token
}

func NewBinary(left Expr, right Expr, operator scanner.Token) Binary {
	return Binary{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (b Binary) Accept(visitor VisitExpr) any {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expr Expr
}

func NewGrouping(expr Expr) Grouping {
	return Grouping{
		Expr: expr,
	}
}

func (g Grouping) Accept(visitor VisitExpr) any {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value any
}

func NewLiteral(value any) Literal {
	return Literal{
		Value: value,
	}
}

func (l Literal) Accept(visitor VisitExpr) any {
	return visitor.VisitLiteralExpr(l)
}

type Logical struct {
	Left     Expr
	Right    Expr
	Operator scanner.Token
}

func NewLogical(left Expr, right Expr, operator scanner.Token) Logical {
	return Logical{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (l Logical) Accept(visitor VisitExpr) any {
	return visitor.VisitLogicalExpr(l)
}

type Unary struct {
	Right    Expr
	Operator scanner.Token
}

func NewUnary(right Expr, operator scanner.Token) Unary {
	return Unary{
		Right:    right,
		Operator: operator,
	}
}

func (u Unary) Accept(visitor VisitExpr) any {
	return visitor.VisitUnaryExpr(u)
}
