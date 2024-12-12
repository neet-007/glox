package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/neet-007/glox/pkg/scanner"
)

/*
 NOTE:
	the timestamp field on the structs is to produce a unique hash for
	each new struct becouse when used in the map for interpeter you need
	to find the right one and dont overwrite
:
*/

type VisitStmt interface {
	VisitClassStmt(stmt Class) (any, error)
	VisitReturnStmt(stmt Return) (any, error)
	VisitFunctionStmt(stmt Function) (any, error)
	VisitVarDeclaration(stmt VarDeclaration) (any, error)
	VisitWhileStmt(stmt WhileStmt) (any, error)
	VisitBlockStmt(stmt Block) (any, error)
	VisitIfStmt(stmt IfStmt) (any, error)
	VisitExpressionStmt(stmt ExpressionStmt) (any, error)
	VisitPrintStmt(stmt PrintStmt) (any, error)
}

type Stmt interface {
	Accept(visitor VisitStmt) (any, error)
}

type Class struct {
	Name       scanner.Token
	Methods    []Function
	SuperClass Variable
	timestamp  int64 // Unique field
}

func NewClass(name scanner.Token, methods []Function, superClass Variable) Class {
	return Class{
		Name:       name,
		Methods:    methods,
		SuperClass: superClass,
		timestamp:  time.Now().UnixNano(),
	}
}

func (c Class) String() string {
	return fmt.Sprintf("class name: %v\nmethods: %v\n\n", c.Name, c.Methods)
}

/*
func (c Class) Equals(other Class) bool {
	if c.Name != other.Name {
		return false
	}

	if len(c.Methods) != len(other.Methods) {
		return false
	}
	for i, value := range c.Methods {
		if other.Methods[i] != value {
			return false
		}
	}

	return true
}
*/

func (c Class) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitClassStmt(c)
}

type Return struct {
	Keyword   scanner.Token
	Value     Expr
	timestamp int64 // Unique field
}

func NewReturn(keyword scanner.Token, value Expr) Return {
	return Return{
		Keyword:   keyword,
		Value:     value,
		timestamp: time.Now().UnixNano(),
	}
}

func (r Return) String() string {
	return fmt.Sprintf("return keyword:\n %v value %v:\n\n", r.Keyword, r.Value)
}

func (r Return) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitReturnStmt(r)
}

type Function struct {
	Name       scanner.Token
	Parameters []scanner.Token
	Body       []Stmt
	timestamp  int64 // Unique field
}

func NewFunction(name scanner.Token, parameters []scanner.Token, body []Stmt) Function {
	return Function{
		Name:       name,
		Parameters: parameters,
		Body:       body,
		timestamp:  time.Now().UnixNano(),
	}
}

func (f Function) String() string {
	return fmt.Sprintf("name:%v\n paramters:%v\n, body:%v\n\n", f.Name, f.Parameters, f.Body)
}

func (f Function) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitFunctionStmt(f)
}

type VarDeclaration struct {
	Initizlier Expr
	Name       scanner.Token
	timestamp  int64 // Unique field
}

func (v VarDeclaration) String() string {
	return fmt.Sprintf("init:%v name:%v\n", v.Initizlier, v.Name)
}

func NewVarDeclaration(name scanner.Token, initizlier Expr) VarDeclaration {
	return VarDeclaration{
		Name:       name,
		Initizlier: initizlier,
		timestamp:  time.Now().UnixNano(),
	}
}

func (v VarDeclaration) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitVarDeclaration(v)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
	timestamp int64 // Unique field
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("body:%v conditno:%v\n", w.Body, w.Condition)
}

func NewWhileStmt(condition Expr, block Stmt) WhileStmt {
	return WhileStmt{
		Condition: condition,
		Body:      block,
		timestamp: time.Now().UnixNano(),
	}
}

func (w WhileStmt) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitWhileStmt(w)
}

type Block struct {
	Statements []Stmt
	timestamp  int64 // Unique field
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
		timestamp:  time.Now().UnixNano(),
	}
}

func (b Block) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitBlockStmt(b)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
	timestamp  int64 // Unique field
}

func (i IfStmt) String() string {
	return fmt.Sprintf("condtion:%v then:%v else:%v\n", i.Condition, i.ThenBranch, i.ElseBranch)
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) IfStmt {
	return IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
		timestamp:  time.Now().UnixNano(),
	}
}

func (i IfStmt) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitIfStmt(i)
}

type ExpressionStmt struct {
	Expression Expr
	timestamp  int64 // Unique field
}

func (e ExpressionStmt) String() string {
	return fmt.Sprintf("expr:%v\n", e.Expression)
}

func NewExpressionStmt(expr Expr) ExpressionStmt {
	return ExpressionStmt{
		Expression: expr,
		timestamp:  time.Now().UnixNano(),
	}
}

func (e ExpressionStmt) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
	timestamp  int64 // Unique field
}

func (p PrintStmt) String() string {
	return fmt.Sprintf("print expr:%v\n", p.Expression)
}

func NewPrintStmt(expr Expr) PrintStmt {
	return PrintStmt{
		Expression: expr,
		timestamp:  time.Now().UnixNano(),
	}
}

func (p PrintStmt) Accept(visitor VisitStmt) (any, error) {
	return visitor.VisitPrintStmt(p)
}
