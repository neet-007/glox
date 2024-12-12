package resolver

import (
	"fmt"

	"github.com/neet-007/glox/pkg/interpreter"
	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/scanner"
)

type FunctionType int
type ClassType int

const (
	NONE_FUNCTION FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

const (
	NONE_CLASS ClassType = iota
	CLASS
	SUBCLASS
)

func (f FunctionType) String() string {
	switch f {
	case NONE_FUNCTION:
		return "NONE_FUNCTION"
	case FUNCTION:
		return "FUNCTION"
	case INITIALIZER:
		return "INITIALIZER"
	case METHOD:
		return "METHOD"
	default:
		return "UNKNOWN_FUNCTION_TYPE"
	}
}

func (c ClassType) String() string {
	switch c {
	case NONE_CLASS:
		return "NONE_CLASS"
	case CLASS:
		return "CLASS"
	case SUBCLASS:
		return "SUBCLASS"
	default:
		return "UNKNOWN_CLASS_TYPE"
	}
}

type CompileError struct {
	Token   scanner.Token
	Message string
}

func (e *CompileError) Error() string {
	return e.Message
}

func NewCompileError(token scanner.Token, message string) *CompileError {
	return &CompileError{
		Token:   token,
		Message: message,
	}
}

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	errors          []*CompileError
	currentFunction FunctionType
	currentClass    ClassType
	debug           bool
}

func NewResolver(interpreter *interpreter.Interpreter, debug bool) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		errors:          []*CompileError{},
		currentFunction: NONE_FUNCTION,
		currentClass:    NONE_CLASS,
		debug:           debug,
	}
}

func (r *Resolver) Resolve(stmts []parser.Stmt) []*CompileError {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}

	return r.errors
}

func (r *Resolver) resolveStmts(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt parser.Stmt) (any, error) {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr parser.Expr) (any, error) {
	return expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr parser.Expr, name scanner.Token) {
	if r.debug {
		fmt.Printf("resolve local name:%s\n", name.Lexeme)
	}
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			if r.debug {
				fmt.Printf("resolve local name:%s found dist %d\n", name.Lexeme, len(r.scopes)-1-i)
			}
			r.interpreter.ResolveExpr(expr, len(r.scopes)-1-i)
			return
		}
	}
	if r.debug {
		fmt.Printf("resolve local name:%s not found\n", name.Lexeme)
	}
}

func (r *Resolver) resolveFunction(stmt parser.Function, functionType FunctionType) {
	if r.debug {
		fmt.Printf("resolve function name: %s type:%s\n", stmt.Name.Lexeme, functionType)
	}
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType
	r.beginScope()

	for _, param := range stmt.Parameters {
		r.declare(param)
		r.define(param)
	}

	if r.debug {
		fmt.Printf("resolve function name: %s type:%s finish\n", stmt.Name.Lexeme, functionType)
	}
	r.resolveStmts(stmt.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) declare(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		r.error(NewCompileError(name, "Already a variable with this name in this scope"))
		return
	}
	scope[name.Lexeme] = false
	return
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

func (r *Resolver) error(errros ...error) {
	for _, err := range errros {
		if compileErr, ok := err.(*CompileError); ok {
			r.errors = append(r.errors, compileErr)
		} else {
			panic("not compile error" + err.Error())
		}
	}
}

func (r *Resolver) VisitSuperExpr(expr parser.Super) (any, error) {
	if r.currentClass == NONE_CLASS {
		r.error(NewCompileError(expr.Keyword, "Can't use 'super' outside of a class"))
		return nil, nil
	} else if r.currentClass != SUBCLASS {
		r.error(NewCompileError(expr.Keyword, "Can't use 'super' in a class with no superclass"))
		return nil, nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr parser.This) (any, error) {
	if r.currentClass == NONE_CLASS {
		r.error(NewCompileError(expr.Keyword, "Can't use 'this' outside of a class"))
		return nil, nil
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr parser.Set) (any, error) {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr parser.Get) (any, error) {
	r.resolveExpr(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr parser.Call) (any, error) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr parser.Variable) (any, error) {
	if len(r.scopes) > 0 {
		if val, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && !val {
			r.error(NewCompileError(expr.Name, "Can't read local variable in its own initializer"))
			return nil, nil
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr parser.Assign) (any, error) {
	r.resolveExpr(expr.Expr)
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

func (r *Resolver) VisitClassStmt(stmt parser.Class) (any, error) {
	if r.debug {
		fmt.Printf("resolver visit class name:%s\n", stmt.Name.Lexeme)
	}
	currentClass := r.currentClass
	r.currentClass = CLASS
	r.declare(stmt.Name)
	r.define(stmt.Name)

	var zeroVariabe parser.Variable
	if stmt.SuperClass != zeroVariabe {
		if r.debug {
			fmt.Printf("resolver visit class name:%s has superclass\n", stmt.Name.Lexeme)
		}
		r.currentClass = SUBCLASS
		if stmt.SuperClass.Name == stmt.Name {
			r.error(NewCompileError(stmt.Name, "A class can't inherit from itself."))
			return nil, nil
		}
		r.resolveExpr(stmt.SuperClass)
	}

	if stmt.SuperClass != zeroVariabe {
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaation := METHOD
		if method.Name.Lexeme == "init" {
			if r.debug {
				fmt.Printf("resolver visit class name:%s has init method\n", stmt.Name.Lexeme)
			}
			declaation = INITIALIZER
		}
		r.resolveFunction(method, declaation)
	}

	r.endScope()

	if stmt.SuperClass != zeroVariabe {
		r.currentClass = currentClass
		r.endScope()
	}

	if r.debug {
		fmt.Printf("resolver visit class name:%s finished\n", stmt.Name.Lexeme)
	}
	r.currentClass = currentClass
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt parser.Return) (any, error) {
	if r.debug {
		fmt.Printf("resolver visit return\n")
	}
	if r.currentFunction == NONE_FUNCTION {
		if r.debug {
			fmt.Printf("resolver visit return not function\n")
		}
		r.error(NewCompileError(stmt.Keyword, "Can't return from top-level code."))
		return nil, nil
	}
	if stmt.Value != nil {
		if r.debug {
			fmt.Printf("resolver visit return has value\n")
		}
		if r.currentFunction == INITIALIZER {
			if r.debug {
				fmt.Printf("resolver visit return has value but in init method\n")
			}
			r.error(NewCompileError(stmt.Keyword, "Can't return a value from an initializer."))
			return nil, nil
		}
		r.resolveExpr(stmt.Value)
	}

	if r.debug {
		fmt.Printf("resolver visit return finished\n")
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
	r.resolveStmts(stmt.Statements)
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
