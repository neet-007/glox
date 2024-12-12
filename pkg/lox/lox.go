package lox

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/neet-007/glox/pkg/interpreter"
	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/resolver"
	"github.com/neet-007/glox/pkg/scanner"
	"github.com/neet-007/glox/pkg/utils"
)

type Lox struct {
	interpreter     *interpreter.Interpreter
	hadError        bool
	hadRuntimeError bool
	debug           bool
	printAst        bool
}

func NewLox() *Lox {
	return &Lox{
		interpreter: interpreter.NewInterpreter(false),
	}
}

func (l *Lox) Main() {
	debug := flag.Bool("debug", false, "turn on debug mode")
	printAst := flag.Bool("ast", false, "print parser AST")
	flag.Parse()

	l.debug = *debug
	l.printAst = *printAst
	l.interpreter.Debug = *debug

	args := flag.Args()

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
	}

	if len(args) == 1 {
		l.runFile(args[0])
	} else {
		l.runPromt()
	}
}

func (l *Lox) runFile(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file %s with error: %w\n", filePath, err)
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
			fmt.Fprintf(os.Stderr, "Error reading input: %w\n", err)
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
	scanner := scanner.NewScanner(source, l.debug)
	tokens, scannerErrors := scanner.Scan()

	for _, err := range scannerErrors {
		l.error(err.Token, err.Message)
	}

	parser_ := parser.NewParser(tokens, l.debug)
	statements, parserErrors := parser_.Parse()

	for _, err := range parserErrors {
		l.error(err.Token, err.Message)
	}

	if l.printAst {
		astPrinter := utils.NewAstPrinter()
		astPrinter.Print(statements)
	}

	if l.hadError {
		return
	}

	resolver_ := resolver.NewResolver(l.interpreter, l.debug)

	compileErros := resolver_.Resolve(statements)
	for _, err := range compileErros {
		l.error(err.Token, err.Message)
	}

	if l.hadError {
		return
	}

	err := l.interpreter.Interpret(statements)
	if err != nil {
		l.error(err.Token, err.Message)
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
