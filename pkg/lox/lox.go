package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/neet-007/glox/pkg/interpreter"
	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/resolver"
	"github.com/neet-007/glox/pkg/scanner"
)

type Lox struct {
	interpreter     *interpreter.Interpreter
	hadError        bool
	hadRuntimeError bool
}

func NewLox() *Lox {
	return &Lox{
		interpreter: interpreter.NewInterpreter(),
	}
}

func (l *Lox) Main() {
	args := os.Args

	if len(args) == 1 {
		l.runPromt()
		return

	}
	if len(args) == 2 {
		l.runFile(args[1])
		return
	}

	fmt.Fprintln(os.Stderr, "Usage: glox [script]")
	os.Exit(64)
}

func (l *Lox) runFile(filePath string) {
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open file with error: %w", err)
	}

	l.run(file)
	if l.hadError {
		os.Exit(65)
	}

	if l.hadRuntimeError {
		os.Exit(70)
	}
}

func (l *Lox) runPromt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			break
		}

		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		l.run(line)
		l.hadError = false
	}
}

func (l *Lox) run(source []byte) {
	scanner := scanner.NewScanner(source)
	tokens, scannerErrors := scanner.Scan()

	for _, err := range scannerErrors {
		l.error(err.Token, err.Message)
	}

	parser_ := parser.NewParser(tokens)
	statements, parserErrors := parser_.Parse()

	for _, err := range parserErrors {
		l.error(err.Token, err.Message)
	}

	if l.hadError {
		return
	}

	interpreter_ := interpreter.NewInterpreter()
	resolver_ := resolver.NewResolver(interpreter_)

	compileErros := resolver_.Resolve(statements)
	for _, err := range compileErros {
		l.error(err.Token, err.Message)
	}

	if l.hadError {
		return
	}

	err := interpreter_.Interpret(statements)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}

}

func (l *Lox) error(token scanner.Token, message string) {
	if token.TokenType == scanner.EOF {
		l.report(token.Line, " at end", message)
	} else {
		l.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	l.hadError = true
}
