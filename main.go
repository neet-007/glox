package main

import (
	"fmt"
	"os"

	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/scanner"
	"github.com/neet-007/glox/pkg/utils"
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
	expressions := parser_.Parse()

	printer := utils.NewAstPrinter()

	for _, expression := range expressions {
		if exprBi, ok := expression.(parser.Binary); ok {
			fmt.Printf("%+v\n", printer.VisitBinaryExpr(exprBi))
		}
		if exprLog, ok := expression.(parser.Logical); ok {
			fmt.Printf("%+v\n", printer.VisitLogicalExpr(exprLog))
			printer.VisitLogicalExpr(exprLog)
		}
		if exprLi, ok := expression.(parser.Literal); ok {
			fmt.Printf("%+v\n", printer.VisitLiteralExpr(exprLi))
		}
		if exprUn, ok := expression.(parser.Unary); ok {
			fmt.Printf("%+v\n", printer.VisitUnaryExpr(exprUn))
		}
		if exprGr, ok := expression.(parser.Grouping); ok {
			fmt.Printf("%+v\n", printer.VisitGroupingExpr(exprGr))
		}
	}

}
