package stateflow

import "fmt"

type TokenType string

const (
	EOF         TokenType = "EOF"
	IDENTIFIER  TokenType = "IDENTIFIER"
	STRING      TokenType = "STRING"
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	ARROW_LEFT  TokenType = "ARROW_LEFT"
	ARROW_RIGHT TokenType = "ARROW_RIGHT"
	NEWLINE     TokenType = "NEWLINE"
	DFA         TokenType = "DFA"
	INITIAL     TokenType = "INITIAL"
	STATE       TokenType = "STATE"
	FINAL       TokenType = "FINAL"
	ON          TokenType = "ON"
	WHEN        TokenType = "WHEN"
	OR          TokenType = "OR"
	FUNCTION    TokenType = "FUNCTION"
	STR         TokenType = "STR"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("(%v %q %s)", t.tokenType, t.lexeme, stringifyTokenValue(t.literal))
}

func stringifyTokenValue(value any) string {
	switch value.(type) {
	case nil:
		return "null"
	}
	return fmt.Sprint(value)
}
