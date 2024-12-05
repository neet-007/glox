package scanner

import (
	"fmt"
	"strconv"
)

type ParseError struct {
	message string
	line    int
}

func newParseError(line int, message string) *ParseError {
	return &ParseError{
		line:    line,
		message: message,
	}
}

func (err *ParseError) Error() string {
	return fmt.Sprintf("parse error at %d with message %s", err.line, err.message)
}

type Scanner struct {
	keywords map[string]TokenType
	tokens   []Token
	source   []byte
	start    int
	current  int
	length   int
	line     int
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{
		keywords: map[string]TokenType{
			"and":    AND,
			"class":  CLASS,
			"else":   ELSE,
			"false":  FALSE,
			"for":    FOR,
			"fun":    FUN,
			"if":     IF,
			"nil":    NIL,
			"or":     OR,
			"print":  PRINT,
			"return": RETURN,
			"super":  SUPER,
			"this":   THIS,
			"true":   TRUE,
			"var":    VAR,
			"while":  WHILE,
		},
		source: source,
		length: len(source),
	}
}

func (s *Scanner) Scan() ([]Token, error) {

	var err error
	for !s.isAtEnd() {
		s.start = s.current
		parseErr := s.scanToken()
		if err != nil {
			err = parseErr
		}
	}

	s.addToken(EOF, nil)
	return s.tokens, err
}

func (s *Scanner) scanToken() error {
	c := s.advance()

	switch c {
	case '(':
		{
			s.addToken(LEFT_PAREN, nil)
			break
		}
	case ')':
		{
			s.addToken(RIGHT_PAREN, nil)
			break

		}
	case '{':
		{
			s.addToken(LEFT_BRACE, nil)
			break

		}
	case '}':
		{
			s.addToken(RIGHT_BRACE, nil)
			break

		}
	case ',':
		{
			s.addToken(COMMA, nil)
		}
	case '.':
		{

			s.addToken(DOT, nil)
		}
	case ';':
		{
			s.addToken(SEMICOLON, nil)
		}
	case '+':
		{
			s.addToken(PLUS, nil)
			break

		}
	case '-':
		{
			s.addToken(MINUS, nil)
			break

		}
	case '*':
		{
			s.addToken(STAR, nil)
			break

		}
	case '/':
		{
			if s.match('/') {
				for !s.isAtEnd() && s.peek() != '\n' {
					s.advance()
				}
				break
			}
			s.addToken(SLASH, nil)
			break

		}
	case '=':
		{
			if s.match('=') {
				s.addToken(EQUAL_EQUAL, nil)
				break
			}
			s.addToken(EQUAL, nil)
			break

		}
	case '>':
		{
			if s.match('=') {
				s.addToken(GREATER_EQUAL, nil)
				break
			}
			s.addToken(GREATER, nil)
			break
		}
	case '<':
		{
			if s.match('=') {
				s.addToken(LESS_EQUAL, nil)
				break
			}
			s.addToken(LESS, nil)
			break
		}
	case '"':
		{
			s.stringLiteral()
			break
		}
	case ' ':
	case '\t':
		{
			break

		}
	case '\n':
		{
			s.line++
			break
		}
	default:
		{
			if s.isNumber(c) {
				s.number()
				break
			}

			if s.isAlphaNumerical(c) {
				s.identifier()
				break
			}

			return newParseError(s.line, "unknown charecter")
		}
	}

	return nil
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	s.tokens = append(s.tokens, Token{
		TokenType: tokenType,
		Literal:   literal,
		Line:      s.line,
		Lexem:     string(s.source[s.start:s.current]),
	})
}

func (s *Scanner) identifier() {
	for !s.isAtEnd() && s.isAlphaNumerical(s.peek()) {
		s.advance()
	}

	keyword, ok := s.keywords[string(s.source[s.start:s.current])]
	if ok {
		s.addToken(keyword, nil)
		return
	}

	s.addToken(IDENTIFIER, string(s.source[s.start:s.current]))
}

func (s *Scanner) stringLiteral() error {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return newParseError(s.line, "unterminated string")
	}

	s.advance()
	s.addToken(STRING, string(s.source[s.start+1:s.current-1]))
	return nil
}

func (s *Scanner) number() error {
	for !s.isAtEnd() && s.isNumber(s.peek()) {
		s.advance()
	}

	if s.match('.') {
		for !s.isAtEnd() && s.isNumber(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.Atoi(string(s.source[s.start:s.current]))
	if err != nil {
		return newParseError(s.line, "invalid number")
	}

	s.addToken(NUMBER, num)
	return nil
}

func (s *Scanner) isNumber(c byte) bool {
	return '0' <= c && c <= '9'
}

func (s *Scanner) isAlphaNumerical(c byte) bool {
	return s.isNumber(c) || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func (s *Scanner) advance() byte {
	returnVal := s.source[s.current]
	s.current++

	return returnVal
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= s.length
}

func (s *Scanner) peek() byte {
	return s.source[s.current]
}

func (s *Scanner) match(c byte) bool {
	if s.peek() == c {
		s.current++
		return true
	}
	return false
}
