package main

import (
	"fmt"
	"os"

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

	fmt.Printf("%v\n", tokens)
}
