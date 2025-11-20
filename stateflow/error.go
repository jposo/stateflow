package stateflow

import (
	"fmt"
)

func report(line int, where string, message string) string {
	return fmt.Sprintf("[line %d] Error%s: %s\n", line, where, message)
}

type SyntaxError struct {
	Line    int
	Message string
}

type ParseError struct {
	Token   *Token
	Message string
}

type RuntimeError struct {
	Token   *Token
	Message string
}

type ZeroDivisionError struct {
	Line    int
	Message string
}

func (s SyntaxError) Error() string {
	return report(s.Line, "", s.Message)
}

func (p ParseError) Error() string {
	where := fmt.Sprintf(" at '%s'", p.Token.lexeme)
	if p.Token.tokenType == EOF {
		where = " at end"
	}
	return report(p.Token.line, where, p.Message)
}

func (r RuntimeError) Error() string {
	return report(r.Token.line, "", r.Message)
}

func (z ZeroDivisionError) Error() string {
	return report(z.Line, "", z.Message)
}
