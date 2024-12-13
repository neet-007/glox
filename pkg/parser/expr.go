package parser

import (
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

type VisitExpr interface {
	VisitListSet(expr ListSet) (any, error)
	VisitListGet(expr ListGet) (any, error)
	VisitListExpr(expr ListExpr) (any, error)
	VisitSuperExpr(expr Super) (any, error)
	VisitThisExpr(expr This) (any, error)
	VisitSetExpr(expr Set) (any, error)
	VisitGetExpr(expr Get) (any, error)
	VisitCallExpr(expr Call) (any, error)
	VisitVariableExpr(expr Variable) (any, error)
	VisitAssignExpr(expr Assign) (any, error)
	VisitBinaryExpr(expr Binary) (any, error)
	VisitGroupingExpr(expr Grouping) (any, error)
	VisitLiteralExpr(expr Literal) (any, error)
	VisitLogicalExpr(expr Logical) (any, error)
	VisitUnaryExpr(expr Unary) (any, error)
}

type Expr interface {
	Accept(visitor VisitExpr) (any, error)
}

type Super struct {
	Keyword   scanner.Token
	Method    scanner.Token
	timestamp int64 // Unique field
}

func NewSuper(keyword scanner.Token, method scanner.Token) Super {
	return Super{
		Keyword:   keyword,
		Method:    method,
		timestamp: time.Now().UnixNano(),
	}
}

func (s Super) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitSuperExpr(s)
}

type This struct {
	Keyword   scanner.Token
	timestamp int64 // Unique field
}

func NewThis(keyword scanner.Token) This {
	return This{
		Keyword:   keyword,
		timestamp: time.Now().UnixNano(),
	}
}

func (t This) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitThisExpr(t)
}

type Set struct {
	Value     Expr
	Object    Expr
	Name      scanner.Token
	timestamp int64 // Unique field
}

func NewSet(value Expr, object Expr, name scanner.Token) Set {
	return Set{
		Value:     value,
		Object:    object,
		Name:      name,
		timestamp: time.Now().UnixNano(),
	}
}

func (s Set) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitSetExpr(s)
}

type Get struct {
	Object    Expr
	Name      scanner.Token
	timestamp int64 // Unique field
}

func NewGet(object Expr, name scanner.Token) Get {
	return Get{
		Object:    object,
		Name:      name,
		timestamp: time.Now().UnixNano(),
	}
}

func (g Get) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitGetExpr(g)
}

type Call struct {
	Callee    Expr
	Paren     scanner.Token
	Arguments []Expr
	timestamp int64 // Unique field
}

func NewCall(callee Expr, paren scanner.Token, arguments []Expr) Call {
	return Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
		timestamp: time.Now().UnixNano(),
	}
}

func (c Call) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitCallExpr(c)
}

type Variable struct {
	Name      scanner.Token
	timestamp int64 // Unique field
}

func NewVariable(name scanner.Token) Variable {
	return Variable{
		Name:      name,
		timestamp: time.Now().UnixNano(),
	}
}

func (v Variable) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Lexem     scanner.Token
	Expr      Expr
	timestamp int64 // Unique field
}

func NewAssign(lexem scanner.Token, expr Expr) Assign {
	return Assign{
		Lexem:     lexem,
		Expr:      expr,
		timestamp: time.Now().UnixNano(),
	}
}

func (a Assign) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitAssignExpr(a)
}

type Binary struct {
	Left      Expr
	Right     Expr
	Operator  scanner.Token
	timestamp int64 // Unique field
}

func NewBinary(left Expr, right Expr, operator scanner.Token) Binary {
	return Binary{
		Left:      left,
		Right:     right,
		Operator:  operator,
		timestamp: time.Now().UnixNano(),
	}
}

func (b Binary) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expr      Expr
	timestamp int64 // Unique field
}

func NewGrouping(expr Expr) Grouping {
	return Grouping{
		Expr:      expr,
		timestamp: time.Now().UnixNano(),
	}
}

func (g Grouping) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value     any
	timestamp int64 // Unique field
}

func NewLiteral(value any) Literal {
	return Literal{
		Value:     value,
		timestamp: time.Now().UnixNano(),
	}
}

func (l Literal) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type ListSet struct {
	List      Expr
	Index     Expr
	Value     Expr
	Token     scanner.Token
	timestamp int64 // Unique field
}

func NewListSet(list Expr, index Expr, value Expr, token scanner.Token) ListSet {
	return ListSet{
		List:      list,
		Index:     index,
		Value:     value,
		Token:     token,
		timestamp: time.Now().UnixNano(),
	}
}

func (l ListSet) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitListSet(l)
}

type ListGet struct {
	List      Expr
	Index     Expr
	Token     scanner.Token
	timestamp int64 // Unique field
}

func NewListGet(list Expr, index Expr, token scanner.Token) ListGet {
	return ListGet{
		List:      list,
		Index:     index,
		Token:     token,
		timestamp: time.Now().UnixNano(),
	}
}

func (l ListGet) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitListGet(l)
}

type ListExpr struct {
	Literals     []Expr
	LeftBracket  scanner.Token
	RightBracket scanner.Token
	timestamp    int64 // Unique field
}

func NewListExpr(leftBracket scanner.Token, rightBracket scanner.Token, literals []Expr) ListExpr {
	return ListExpr{
		LeftBracket:  leftBracket,
		RightBracket: rightBracket,
		Literals:     literals,
		timestamp:    time.Now().UnixNano(),
	}
}

func (l ListExpr) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitListExpr(l)
}

type Logical struct {
	Left      Expr
	Right     Expr
	Operator  scanner.Token
	timestamp int64 // Unique field
}

func NewLogical(left Expr, right Expr, operator scanner.Token) Logical {
	return Logical{
		Left:      left,
		Right:     right,
		Operator:  operator,
		timestamp: time.Now().UnixNano(),
	}
}

func (l Logical) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

type Unary struct {
	Right     Expr
	Operator  scanner.Token
	timestamp int64 // Unique field
}

func NewUnary(right Expr, operator scanner.Token) Unary {
	return Unary{
		Right:     right,
		Operator:  operator,
		timestamp: time.Now().UnixNano(),
	}
}

func (u Unary) Accept(visitor VisitExpr) (any, error) {
	return visitor.VisitUnaryExpr(u)
}
