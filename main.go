package main

import (
	"fmt"
	"os"

	"github.com/jposo/stateflow/stateflow"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	fileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := stateflow.Scanner{Source: fileContents}
	_, scanErrs := scanner.ScanTokens()
	scanner.PrintTokens()
	if len(scanErrs) > 0 {
		for _, err := range scanErrs {
			fmt.Fprint(os.Stderr, err.Error())
		}
		os.Exit(65) // Syntax Error
	}
}
