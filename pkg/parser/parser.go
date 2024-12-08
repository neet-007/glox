package parser

import (
	"fmt"

	"github.com/neet-007/glox/pkg/scanner"
)

type ParseError struct {
	Token   scanner.Token
	Message string
}

func newParseError(token scanner.Token, message string) *ParseError {
	return &ParseError{
		Token:   token,
		Message: message,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("parse error at %d with message %s\n", p.Token.Line, p.Message)
}

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

func (p *Parser) Parse() ([]Stmt, []*ParseError) {
	expressions := []Stmt{}
	errors := []*ParseError{}

	for !p.isAtEnd() {
		statement, parseErr := p.declaration()
		if parseErr != nil {
			errors = append(errors, parseErr)
			p.synchronize()
		}
		expressions = append(expressions, statement)
	}

	return expressions, errors
}

func (p *Parser) declaration() (Stmt, *ParseError) {
	if p.match(scanner.FUN) {
		return p.function("function")
	}
	if p.match(scanner.FOR) {
		return p.forStatement()
	}
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) function(kind string) (Stmt, *ParseError) {
	name, parseErr := p.consume(scanner.IDENTIFIER, "Expect function "+kind+" name")
	if parseErr != nil {
		return nil, parseErr
	}

	_, parseErr = p.consume(scanner.LEFT_PAREN, "Expect '(' for function")
	if parseErr != nil {
		return nil, parseErr
	}

	parameters := []scanner.Token{}
	var paramSizeErr *ParseError
	if !p.check(scanner.RIGHT_PAREN) {
		_, parseErr = p.consume(scanner.IDENTIFIER, "Expect identefier for parameter")
		if parseErr != nil {
			return nil, parseErr
		}

		parameters = append(parameters, p.previous())
		for p.match(scanner.COMMA) {
			if len(parameters) >= 255 {
				paramSizeErr = newParseError(name, kind+"s have a max of 256 parameters")
			}
			_, parseErr = p.consume(scanner.IDENTIFIER, "Expect identefier for parameter")
			if parseErr != nil {
				return nil, parseErr
			}
			parameters = append(parameters, p.previous())
		}
	}

	_, parseErr = p.consume(scanner.RIGHT_PAREN, "Expect ')' for function")
	if parseErr != nil {
		return nil, parseErr
	}

	_, parseErr = p.consume(scanner.LEFT_BRACE, "Expect '{' for block")
	if parseErr != nil {
		return nil, parseErr
	}

	body, parseErr := p.block()
	if parseErr != nil {
		return nil, parseErr
	}

	if paramSizeErr != nil {
		return nil, paramSizeErr
	}

	return NewFunction(name, parameters, body), nil
}

func (p *Parser) forStatement() (Stmt, *ParseError) {
	_, parseErr := p.consume(scanner.LEFT_PAREN, "Expect '(' after for statement")
	if parseErr != nil {
		return nil, parseErr
	}

	var initizlier Stmt
	if p.match(scanner.SEMICOLON) {
		initizlier = nil
	} else if p.match(scanner.VAR) {
		initizlier, parseErr = p.varDeclaration()
		if parseErr != nil {
			return nil, parseErr
		}
	} else {
		initizlier, parseErr = p.expressionStatement()
		if parseErr != nil {
			return nil, parseErr
		}
	}

	var condition Expr
	if !p.check(scanner.SEMICOLON) {
		condition, parseErr = p.expression()
		if parseErr != nil {
			return nil, parseErr
		}
	}
	_, parseErr = p.consume(scanner.SEMICOLON, "Expect ';' after condition")
	if parseErr != nil {
		return nil, parseErr
	}

	var increment Expr
	if !p.check(scanner.RIGHT_PAREN) {
		increment, parseErr = p.expression()
		if parseErr != nil {
			return nil, parseErr
		}
	}
	_, parseErr = p.consume(scanner.RIGHT_PAREN, "Expect ')' after increment")
	if parseErr != nil {
		return nil, parseErr
	}

	body, parseErr := p.statement()
	if parseErr != nil {
		return nil, parseErr
	}

	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpressionStmt(increment)})
	}

	if condition == nil {
		condition = NewLiteral(true)
	}

	body = NewWhileStmt(condition, body)

	if initizlier != nil {
		body = NewBlock([]Stmt{initizlier, body})
	}

	return body, nil
}

func (p *Parser) varDeclaration() (Stmt, *ParseError) {
	identifier, parserErr := p.consume(scanner.IDENTIFIER, "Expect identefier for variable")
	if parserErr != nil {
		return nil, parserErr
	}

	var initilizer Expr
	if p.match(scanner.EQUAL) {
		initilizer, parserErr = p.expression()
		if parserErr != nil {
			return nil, parserErr
		}
	}

	_, parserErr = p.consume(scanner.SEMICOLON, "Expect ';' after expression")
	return NewVarDeclaration(identifier, initilizer), nil
}

func (p *Parser) statement() (Stmt, *ParseError) {
	if p.match(scanner.PRINT) {
		expr, parseErr := p.expression()
		if parseErr != nil {
			return nil, parseErr
		}

		_, parseErr = p.consume(scanner.SEMICOLON, "Expect ';' after expression")
		if parseErr != nil {
			return nil, parseErr
		}
		return NewPrintStmt(expr), nil
	}
	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.RETURN) {
		return p.returnStatemnt()
	}
	if p.match(scanner.LEFT_BRACE) {
		statements, parseErr := p.block()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewBlock(statements), nil
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatemnt() (Stmt, *ParseError) {
	keyword := p.previous()
	var val Expr = nil
	var parseErr *ParseError
	if !p.check(scanner.SEMICOLON) {
		val, parseErr = p.expression()
		if parseErr != nil {
			return nil, parseErr
		}
	}

	_, parseErr = p.consume(scanner.SEMICOLON, "Expect ';' after expression")
	if parseErr != nil {
		return nil, parseErr
	}

	return NewReturn(keyword, val), nil
}

func (p *Parser) whileStatement() (Stmt, *ParseError) {
	_, parseErr := p.consume(scanner.LEFT_PAREN, "Expect '(' afer if statemnt")
	if parseErr != nil {
		return nil, parseErr
	}
	expr, parseErr := p.expression()

	if parseErr != nil {
		return nil, parseErr
	}
	_, parseErr = p.consume(scanner.RIGHT_PAREN, "Expect ')' afer if statemnt")
	if parseErr != nil {
		return nil, parseErr
	}

	body, parseErr := p.statement()
	if parseErr != nil {
		return nil, parseErr
	}

	return NewWhileStmt(expr, body), nil
}

func (p *Parser) ifStatement() (Stmt, *ParseError) {
	_, parseErr := p.consume(scanner.LEFT_PAREN, "Expect '(' afer if statemnt")
	if parseErr != nil {
		return nil, parseErr
	}

	expr, parseErr := p.expression()
	if parseErr != nil {
		return nil, parseErr
	}

	_, parseErr = p.consume(scanner.RIGHT_PAREN, "Expect ')' afer if statemnt")
	if parseErr != nil {
		return nil, parseErr
	}

	_, parseErr = p.consume(scanner.LEFT_BRACE, "Expect '{' block start")
	if parseErr != nil {
		return nil, parseErr
	}

	ifBracnh, parseErr := p.statement()
	if parseErr != nil {
		return nil, parseErr
	}

	var elseBranch Stmt
	if p.match(scanner.ELSE) {
		_, parseErr = p.consume(scanner.LEFT_BRACE, "Expect '{' block after if statemnt")
		if parseErr != nil {
			return nil, parseErr
		}

		elseBranch, parseErr = p.statement()
		if parseErr != nil {
			return nil, parseErr
		}

	}
	return NewIfStmt(expr, ifBracnh, elseBranch), nil
}

func (p *Parser) block() ([]Stmt, *ParseError) {
	statemnts := []Stmt{}
	for !p.isAtEnd() && !p.check(scanner.RIGHT_BRACE) {
		statement, parseErr := p.declaration()
		if parseErr != nil {
			return nil, parseErr
		}

		statemnts = append(statemnts, statement)
	}

	_, parseErr := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block")
	if parseErr != nil {
		return nil, parseErr
	}

	return statemnts, nil
}

func (p *Parser) expressionStatement() (Stmt, *ParseError) {
	expr, parseErr := p.expression()
	if parseErr != nil {
		return nil, parseErr
	}

	_, parseErr = p.consume(scanner.SEMICOLON, "Expect ';' after expression")
	if parseErr != nil {
		return nil, parseErr
	}

	return NewExpressionStmt(expr), nil
}

func (p *Parser) expression() (Expr, *ParseError) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, *ParseError) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		equal := p.previous()
		val, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if exprVal, ok := expr.(Variable); ok {
			return NewAssign(exprVal.Name, val), nil
		}

		_, err = p.consume(scanner.SEMICOLON, "Expect ';' after expression")
		if err != nil {
			return nil, err
		}
		return nil, newParseError(equal, "assigenmnt to invalid value")
	}

	return expr, nil
}

func (p *Parser) or() (Expr, *ParseError) {
	left, parseErr := p.and()
	if parseErr != nil {
		return nil, parseErr
	}

	if p.match(scanner.OR) {
		operator := p.previous()
		right, parseErr := p.or()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewLogical(left, right, operator), nil

	}

	return left, nil
}

func (p *Parser) and() (Expr, *ParseError) {
	left, parseErr := p.equality()
	if parseErr != nil {
		return nil, parseErr
	}

	if p.match(scanner.AND) {
		operator := p.previous()
		right, parseErr := p.and()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewLogical(left, right, operator), nil

	}

	return left, nil
}

func (p *Parser) equality() (Expr, *ParseError) {
	left, parseErr := p.comparison()
	if parseErr != nil {
		return nil, parseErr
	}

	if p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right, parseErr := p.equality()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewLogical(left, right, operator), nil
	}

	return left, nil
}

func (p *Parser) comparison() (Expr, *ParseError) {
	left, parseErr := p.term()
	if parseErr != nil {
		return nil, parseErr
	}

	if p.match(scanner.GREATER, scanner.GREATER, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right, parseErr := p.comparison()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewLogical(left, right, operator), nil
	}

	return left, nil
}

func (p *Parser) term() (Expr, *ParseError) {
	left, parseErr := p.factor()
	if parseErr != nil {
		return nil, parseErr
	}

	if p.match(scanner.PLUS, scanner.MINUS) {
		operator := p.previous()
		rigth, parseErr := p.term()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewBinary(left, rigth, operator), nil
	}

	return left, nil
}

func (p *Parser) factor() (Expr, *ParseError) {
	left, parseErr := p.unary()
	if parseErr != nil {
		return nil, parseErr
	}

	if p.match(scanner.STAR, scanner.SLASH) {
		operator := p.previous()
		rigth, parseErr := p.factor()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewBinary(left, rigth, operator), nil
	}

	return left, nil
}

func (p *Parser) unary() (Expr, *ParseError) {
	if p.match(scanner.MINUS, scanner.BANG) {
		operator := p.previous()
		right, parseErr := p.unary()
		if parseErr != nil {
			return nil, parseErr
		}

		return NewUnary(right, operator), nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, *ParseError) {
	name, parseErr := p.primary()
	if parseErr != nil {
		return nil, parseErr
	}

	for {
		if p.match(scanner.LEFT_PAREN) {
			name, parseErr = p.finishCall(name)
		} else {
			break
		}
	}

	return name, nil
}

func (p *Parser) finishCall(expr Expr) (Expr, *ParseError) {
	arguments := []Expr{}

	var argumentSizeErr *ParseError
	if !p.check(scanner.RIGHT_PAREN) {
		expr, parseErr := p.expression()
		if parseErr != nil {
			return nil, parseErr
		}

		arguments = append(arguments, expr)

		for p.match(scanner.COMMA) {
			if len(arguments) > 255 {
				argumentSizeErr = newParseError(scanner.Token{}, "calls have a max of 256 parameters")
			}

			expr, parseErr := p.expression()
			if parseErr != nil {
				return nil, parseErr
			}

			arguments = append(arguments, expr)
		}
	}

	paren, parseErr := p.consume(scanner.RIGHT_PAREN, "Expect ')' after call")
	if parseErr != nil {
		return nil, parseErr
	}

	if argumentSizeErr != nil {
		argumentSizeErr.Token = paren
		return nil, argumentSizeErr
	}
	return NewCall(expr, paren, arguments), nil
}

func (p *Parser) primary() (Expr, *ParseError) {
	if p.match(scanner.TRUE) {
		return NewLiteral(true), nil
	}
	if p.match(scanner.FALSE) {
		return NewLiteral(false), nil
	}
	if p.match(scanner.NIL) {
		return NewLiteral(nil), nil
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return NewLiteral(p.previous().Literal), nil
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, parseErr := p.expression()
		if parseErr != nil {
			return nil, parseErr
		}

		_, parseErr = p.consume(scanner.RIGHT_PAREN, "Expect ')' after grouping")

		if parseErr != nil {
			return nil, parseErr
		}

		return NewGrouping(expr), nil
	}
	if p.match(scanner.IDENTIFIER) {
		return NewVariable(p.previous()), nil
	}

	return nil, newParseError(p.peek(), "invalid primary")
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, *ParseError) {
	if p.check(tokenType) {
		return p.advnace(), nil
	}

	return scanner.Token{}, newParseError(p.peek(), message)
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

func (p *Parser) synchronize() {
	p.advnace()

	for !p.isAtEnd() {
		if p.previous().TokenType == scanner.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case scanner.CLASS:
		case scanner.FUN:
		case scanner.VAR:
		case scanner.FOR:
		case scanner.IF:
		case scanner.WHILE:
		case scanner.PRINT:
		case scanner.RETURN:
			return
		}

		p.advnace()
	}
}
