package stateflow

import "fmt"

type TokenType string

const (
	EOF            TokenType = "EOF"
	IDENTIFIER     TokenType = "IDENTIFIER"
	STRING_LITERAL TokenType = "STRING_LITERAL"
	LEFT_PAREN     TokenType = "LEFT_PAREN"
	RIGHT_PAREN    TokenType = "RIGHT_PAREN"
	LEFT_BRACE     TokenType = "LEFT_BRACE"
	RIGHT_BRACE    TokenType = "RIGHT_BRACE"
	ARROW_LEFT     TokenType = "ARROW_LEFT"
	ARROW_RIGHT    TokenType = "ARROW_RIGHT"
	NEWLINE        TokenType = "NEWLINE"
	SEMICOLON      TokenType = "SEMICOLON"
	DFA            TokenType = "DFA"
	INITIAL        TokenType = "INITIAL"
	STATE          TokenType = "STATE"
	FINAL          TokenType = "FINAL"
	ON             TokenType = "ON"
	WHEN           TokenType = "WHEN"
	OR             TokenType = "OR"
	FUNCTION       TokenType = "FUNCTION"
	STRING         TokenType = "STRING"
	REGEX          TokenType = "REGEX"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("(%v %q %s)", t.tokenType, t.lexeme, t.stringifyLiteral())
}

func (t Token) stringifyLiteral() string {
	switch t.literal.(type) {
	case nil:
		return "null"
	}
	return fmt.Sprint(t.literal)
}
