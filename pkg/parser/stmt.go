package parser

import "github.com/neet-007/glox/pkg/scanner"

type VisitStmt interface {
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

type VarDeclaration struct {
	Initizlier Expr
	Name       scanner.Token
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

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) IfStmt {
	return IfStmt{condition, thenBranch, elseBranch}
}

func (i IfStmt) Accept(visitor VisitStmt) any {
	return visitor.VisitIfStmt(i)
}

type ExpressionStmt struct {
	Expression Expr
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

func NewPrintStmt(expr Expr) PrintStmt {
	return PrintStmt{
		Expression: expr,
	}
}

func (p PrintStmt) Accept(visitor VisitStmt) any {
	return visitor.VisitPrintStmt(p)
}
