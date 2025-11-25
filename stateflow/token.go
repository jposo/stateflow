package stateflow

import (
	"fmt"
	"math"
)

type TokenType string

const (
	BOF            TokenType = "BOF"
	EOF            TokenType = "EOF"
	IDENTIFIER     TokenType = "IDENTIFIER"
	STRING_LITERAL TokenType = "STRING_LITERAL"
	LEFT_PAREN     TokenType = "LEFT_PAREN"
	RIGHT_PAREN    TokenType = "RIGHT_PAREN"
	LEFT_BRACE     TokenType = "LEFT_BRACE"
	RIGHT_BRACE    TokenType = "RIGHT_BRACE"
	ARROW_LEFT     TokenType = "ARROW_LEFT"
	ARROW_RIGHT    TokenType = "ARROW_RIGHT"
	SEMICOLON      TokenType = "DELIMITER"
	COMMA          TokenType = "COMMA"
	DFA            TokenType = "DFA"
	NFA            TokenType = "NFA"
	INITIAL        TokenType = "INITIAL"
	STATE          TokenType = "STATE"
	FINAL          TokenType = "FINAL"
	ON             TokenType = "ON"
	WHEN           TokenType = "WHEN"
	OR             TokenType = "OR"
	FUNCTION       TokenType = "FUNCTION"
	QUOTE          TokenType = "QUOTE"
	SLASH          TokenType = "SLASH"
	STRING         TokenType = "STRING"
	REGEX          TokenType = "REGEX"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("(%v %q Line %d)", t.tokenType, t.lexeme, t.line)
}

// func (t Token) stringifyLiteral() string {
// 	switch t.literal.(type) {
// 	case nil:
// 		return "nil"
// 	}
// 	return fmt.Sprint(t.literal)
// }

func stringifyTokenValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case float64:
		return printFloat64(v)
	}
	return fmt.Sprint(value)
}

func printFloat64(f64 float64) string {
	if f64 == math.Floor(f64) {
		return fmt.Sprintf("%.1f", f64)
	}
	return fmt.Sprintf("%g", f64)
}
