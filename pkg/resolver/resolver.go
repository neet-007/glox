package resolver

import (
	"github.com/neet-007/glox/pkg/interpreter"
	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/scanner"
)

type FunctionType int

const (
	NONE = iota
	FUNCTION
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NONE,
	}
}

func (r *Resolver) ResolveStms(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt parser.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr parser.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr parser.Expr, name scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.ResolveExpr(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(stmt parser.Function, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType
	r.beginScope()

	for _, param := range stmt.Parameters {
		r.declare(param)
		r.define(param)
	}

	r.ResolveStms(stmt.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) declare(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		//!TODO error wtih Already a variable with this name in this scope.
		return
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitCallExpr(expr parser.Call) (any, error) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr parser.Variable) (any, error) {
	if len(r.scopes) > 0 && !r.scopes[len(r.scopes)-1][expr.Name.Lexeme] {
		//!TODO error with Can't read local variable in its own initializer.
		return nil, nil
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr parser.Assign) (any, error) {
	r.resolveExpr(expr)
	r.resolveLocal(expr, expr.Lexem)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr parser.Binary) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr parser.Grouping) (any, error) {
	r.resolveExpr(expr.Expr)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr parser.Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr parser.Logical) (any, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr parser.Unary) (any, error) {
	r.resolveExpr(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt parser.Return) (any, error) {
	if r.currentFunction == NONE {
		//!TODO error with Can't return from top-level code.
		return nil, nil
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt parser.Function) (any, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUNCTION)
	return nil, nil
}

func (r *Resolver) VisitVarDeclaration(stmt parser.VarDeclaration) (any, error) {
	r.declare(stmt.Name)
	if stmt.Initizlier != nil {
		r.resolveExpr(stmt.Initizlier)
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt parser.WhileStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil, nil
}

func (r *Resolver) VisitBlockStmt(stmt parser.Block) (any, error) {
	r.beginScope()
	r.ResolveStms(stmt.Statements)
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt parser.ExpressionStmt) (any, error) {
	r.resolveExpr(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt parser.PrintStmt) (any, error) {
	r.resolveExpr(stmt.Expression)
	return nil, nil
}
