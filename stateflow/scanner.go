package stateflow

import (
	"errors"
	"fmt"
)

var keywords = map[string]TokenType{
	"dfa":     DFA,
	"initial": INITIAL,
	"final":   FINAL,
	"on":      ON,
	"when":    WHEN,
	"or":      OR,
	"fn":      FUNCTION,
	"str":     STR,
}

type Scanner struct {
	Source  []byte
	tokens  []Token
	start   int
	current int
	line    int
	errors  []error
}

func (s *Scanner) PrintTokens() {
	for _, token := range s.tokens {
		fmt.Println(token)
	}
}

func (s *Scanner) ScanTokens() ([]Token, []error) {
	s.start = 0
	s.current = 0
	s.line = 1
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens, s.errors
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case '-':
		if s.match('>') {
			s.addToken(ARROW_RIGHT)
		}
	case '<':
		if s.match('-') {
			s.addToken(ARROW_LEFT)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		}
	case '"':
		err := s.string()
		if err != nil {
			return err
		}
	case '\n':
		prev := s.glance()
		if prev != nil && (prev.tokenType == IDENTIFIER ||
			prev.tokenType == STRING) {
			fmt.Println("Can insert semicolon in line ", s.line)
		}
		s.addToken(NEWLINE)
		s.line += 1
	case ' ':
	case '\t':
	case '\r':
	default:
		if isAlpha(c) {
			s.identifier()
		} else {
			return fmt.Errorf("Unexpected character: %s\n", s.Source[s.start:s.current])
		}
	}
	return nil
}

// Returns the current byte in source and advances to next byte
func (s *Scanner) advance() byte {
	char := s.Source[s.current]
	s.current += 1
	return char
}

// Store token, optional literal
func (s *Scanner) addToken(tokenType TokenType, literal ...any) {
	var lit any = nil
	if len(literal) > 0 {
		lit = literal[0]
	}
	lexeme := string(s.Source[s.start:s.current])
	s.tokens = append(s.tokens, Token{tokenType, lexeme, lit, s.line})
}

// Verifies in next byte is as expected, if it is, advances to next byte
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.current] != expected {
		return false
	}
	s.current += 1
	return true
}

// Look at current byte
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return byte(0)
	}
	return s.Source[s.current]
}

// Look at previous token
func (s *Scanner) glance() *Token {
	if len(s.tokens) == 0 {
		return nil
	}
	return &s.tokens[len(s.tokens)-1]
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}
	if s.isAtEnd() {
		return errors.New("Unterminated string.\n")
	}
	s.advance() // Closing "

	value := string(s.Source[s.start+1 : s.current-1])
	s.addToken(STRING, value)
	return nil
}

func (s *Scanner) identifier() {
	for isAlphanumeric(s.peek()) {
		s.advance()
	}
	text := string(s.Source[s.start:s.current])
	value, ok := keywords[text]
	if !ok {
		value = IDENTIFIER
	}
	s.addToken(value)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphanumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}
