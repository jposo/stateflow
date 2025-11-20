package main

import (
	"fmt"
	"os"

	"github.com/jposo/stateflow/stateflow"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: stateflow tokenize <filename>")
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

	switch op {
	case "tokenize":
		scanner.PrintTokens()
		if len(scanErrs) > 0 {
			for _, err := range scanErrs {
				fmt.Fprint(os.Stderr, err.Error())
			}
			os.Exit(65) // Lexical Error
		}
	case "parse":
		if len(scanErrs) > 0 {
			os.Exit(65)
		}
		parser := stateflow.Parser{Tokens: tokens}
		_, parseErr := parser.Parse()
		if parseErr != nil {
			fmt.Fprint(os.Stderr, parseErr.Error())
			os.Exit(65) // Syntax Error
		}
		fmt.Println("No errors!")
		// printer := ste.AstPrinter{}
		// printer.Print(expression)
	default:
		fmt.Fprintf(os.Stderr, "Invalid operation.")
	}

}
