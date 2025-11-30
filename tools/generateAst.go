package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func defineAst(outputPath, baseName string, types []string) {
	writer := fmt.Sprintf(`
	package stateflow
	type %s interface {`, baseName)
	writer += fmt.Sprintf(`
		Accept(visitor %sVisitor) (any, error)
	}
	`, baseName)

	defineVisitor(&writer, baseName, types)

	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(&writer, baseName, className, fields)
	}
	f, err := os.Create(outputPath)
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Println(cwd)
		panic(err)
	}
	defer f.Close()
	f.WriteString(writer)
	cmd := exec.Command("go", "fmt", outputPath)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Could not format generated file.")
	}
}

func defineVisitor(writer *string, baseName string, types []string) {
	*writer += fmt.Sprintf(`
		type %sVisitor interface {`, baseName)

	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		*writer += fmt.Sprintf(`
		Visit%s(%s %s) (any, error)`, typeName+baseName, strings.ToLower(baseName), typeName)
	}
	*writer += fmt.Sprint(`
	}
	`)
}

func defineType(writer *string, baseName string, className string, fieldList string) {
	lowercaseLetter := byte(unicode.ToLower(rune(className[0])))
	*writer += fmt.Sprintf(`
	type %s struct {`, className)

	fields := strings.Split(fieldList, ", ")
	for _, f := range fields {
		name := strings.Split(f, " ")[0]
		t := strings.Split(f, " ")[1]
		*writer += fmt.Sprintf(`
		%s %s`, name, t)
	}
	*writer += fmt.Sprint(`
	}`)

	*writer += fmt.Sprintf(`
	func (%c %s) Accept (visitor %sVisitor) (any, error) {
		return visitor.Visit%s(%c)
	}`, lowercaseLetter, className, baseName, className+baseName, lowercaseLetter)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(1)
	}

	command := os.Args[1]
	outputDir := os.Args[2]

	switch command {
	case "generate_ast":
		defTypes := []string{
			"AutomatonDef:autType Token, name Token, stmts []Stmt",
			"FunctionDef:name Token, params []Token, statements []Statement",
		}
		defPath := filepath.Join(outputDir, "definition.go")
		defineAst(defPath, "Definition", defTypes)

		statementTypes := []string{
			"Call: target Token, input Token",
		}
		statementPath := filepath.Join(outputDir, "statement.go")
		defineAst(statementPath, "Statement", statementTypes)

		stmtTypes := []string{
			"StateDecl:stateType Token, name Token",
			"TransDecl:symbol Token, fromState Token, toState Token, conditions []Condition",
		}
		stmtPath := filepath.Join(outputDir, "stmt.go")
		defineAst(stmtPath, "Stmt", stmtTypes)

		conditionTypes := []string{
			"StringCondition:value string",
			"RegexCondition:pattern string",
		}
		conditionPath := filepath.Join(outputDir, "condition.go")
		defineAst(conditionPath, "Condition", conditionTypes)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
	os.Exit(0)
}
