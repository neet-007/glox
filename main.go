package main

import (
	"os"

	"github.com/neet-007/glox/pkg/interpreter"
	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/scanner"
)

func main() {
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	scanner := scanner.NewScanner(file)
	tokens, err := scanner.Scan()
	if err != nil {
		panic(err)
	}

	parser_ := parser.NewParser(tokens)
	statements := parser_.Parse()

	interpreter_ := interpreter.NewInterpreter()
	interpreter_.Interpret(statements)

	/*
		printer := utils.NewAstPrinter()

		printer.Print(statements)
	*/

}
