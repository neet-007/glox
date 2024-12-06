package parser

import (
	"github.com/neet-007/glox/pkg/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	length  int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens: tokens,
		length: len(tokens),
	}
}

func (p *Parser) Parse() []Stmt {
	expressions := []Stmt{}

	for !p.isAtEnd() {
		expressions = append(expressions, p.statement())
	}

	return expressions
}

func (p *Parser) statement() Stmt {
	if p.match(scanner.PRINT) {
		expr := p.expression()
		return NewPrintStmt(expr)
	}

	return p.expressionStatement()
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "Expect ';' after expression")
	return NewExpressionStmt(expr)
}

func (p *Parser) expression() Expr {
	return p.or()
}

func (p *Parser) or() Expr {
	left := p.and()

	if p.match(scanner.OR) {
		operator := p.previous()
		right := p.or()

		return NewLogical(left, right, operator)

	}

	return left
}

func (p *Parser) and() Expr {
	left := p.equality()

	if p.match(scanner.AND) {
		operator := p.previous()
		right := p.and()

		return NewLogical(left, right, operator)

	}

	return left
}

func (p *Parser) equality() Expr {
	left := p.comparison()

	if p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.equality()

		return NewLogical(left, right, operator)
	}

	return left
}

func (p *Parser) comparison() Expr {
	left := p.term()

	if p.match(scanner.GREATER, scanner.GREATER, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.comparison()

		return NewLogical(left, right, operator)
	}

	return left
}

func (p *Parser) term() Expr {
	left := p.factor()

	if p.match(scanner.PLUS, scanner.MINUS) {
		operator := p.previous()
		rigth := p.term()

		return NewBinary(left, rigth, operator)
	}

	return left
}

func (p *Parser) factor() Expr {
	left := p.unary()

	if p.match(scanner.STAR, scanner.SLASH) {
		operator := p.previous()
		rigth := p.factor()

		return NewBinary(left, rigth, operator)
	}

	return left
}

func (p *Parser) unary() Expr {
	if p.match(scanner.MINUS, scanner.BANG) {
		operator := p.previous()
		right := p.unary()

		return NewUnary(right, operator)
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.TRUE) {
		return NewLiteral(true)
	}
	if p.match(scanner.FALSE) {
		return NewLiteral(false)
	}
	if p.match(scanner.NIL) {
		return NewLiteral(nil)
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return NewLiteral(p.previous().Literal)
	}
	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHT_PAREN, "Expect ')' after grouping")
		return NewGrouping(expr)
	}

	//!TODO error
	return NewLiteral(false)
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) scanner.Token {
	if p.check(tokenType) {
		return p.advnace()
	}

	//!TIDI error
	return scanner.Token{}
}

func (p *Parser) advnace() scanner.Token {
	p.current++

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == scanner.EOF
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	return p.peek().TokenType == tokenType
}

func (p *Parser) match(tokenTypes ...scanner.TokenType) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advnace()
			return true
		}
	}

	return false
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) peekAhead() scanner.Token {
	if p.current+1 >= p.length {
		//!TODO error
	}
	return p.tokens[p.current+1]
}
