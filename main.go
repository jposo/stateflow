package main

import (
	"fmt"
	"os"

	"github.com/jposo/stateflow/stateflow"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: stateflow (tokenize | parse) <filename>")
		os.Exit(1)
	}
	op := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := stateflow.Scanner{Source: fileContents}
	tokens, scanErrs := scanner.ScanTokens()
	if len(scanErrs) > 0 {
		for _, err := range scanErrs {
			fmt.Fprint(os.Stderr, err.Error())
		}
		os.Exit(65) // Lexical Error
	}

	switch op {
	case "tokenize":
		scanner.PrintTokens()
	case "parse":
		parser := stateflow.Parser{Tokens: tokens}
		_, parseErr := parser.Parse()
		if parseErr != nil {
			fmt.Fprint(os.Stderr, parseErr.Error())
			os.Exit(65) // Syntax or Semantics Error
		}
		fmt.Println("No errors!")
		// for _, def := range defs {
		// 	fmt.Println(def)
		// }
		// printer := stateflow.AstPrinter{}

	default:
		fmt.Fprintf(os.Stderr, "Invalid operation.")
	}
	os.Exit(0)
}
