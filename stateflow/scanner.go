package stateflow

import (
	"errors"
	"fmt"
)

var keywords = map[string]TokenType{
	"dfa":     DFA,
	"state":   STATE,
	"initial": INITIAL,
	"final":   FINAL,
	"on":      ON,
	"when":    WHEN,
	"or":      OR,
	"fn":      FUNCTION,
	"str":     STRING,
}

type Scanner struct {
	Source  []byte
	tokens  []Token
	start   int
	current int
	line    int
	// inString bool
	errors []error
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
	s.tokens = append(s.tokens, Token{BOF, "", s.line})
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			s.errors = append(s.errors, err)
		}
	}
	s.tokens = append(s.tokens, Token{EOF, "", s.line})
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
		} else {
			return s.unexpectedError()
		}
	case '<':
		if s.match('-') {
			s.addToken(ARROW_LEFT)
		} else {
			return s.unexpectedError()
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			err := s.regex()
			if err != nil {
				return err
			}
		}
	case '"':
		// s.inString = !s.inString
		// s.addToken(QUOTE)
		err := s.string()
		if err != nil {
			return err
		}
	case ';':
		s.addToken(SEMICOLON)
	case '\n':
		prev := s.lastToken()
		if prev != nil {
			switch prev.tokenType {
			case IDENTIFIER, STRING_LITERAL:
				s.addToken(SEMICOLON)
			}
		}
		s.line += 1
	case ' ':
	case '\t':
	case '\r':
	default:
		if isAlpha(c) {
			s.identifier()
		} else {
			return s.unexpectedError()
		}
	}
	return nil
}

func (s *Scanner) unexpectedError() error {
	return fmt.Errorf("Unexpected character in line %d: %s\n", s.line, s.Source[s.start:s.current])
}

// Returns the current byte in source and advances to next byte
func (s *Scanner) advance() byte {
	char := s.Source[s.current]
	s.current += 1
	return char
}

// Store token, optional literal
func (s *Scanner) addToken(tokenType TokenType) {
	lexeme := string(s.Source[s.start:s.current])
	s.tokens = append(s.tokens, Token{tokenType, lexeme, s.line})
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
func (s *Scanner) lastToken() *Token {
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

	s.addToken(STRING_LITERAL)
	return nil
}

func (s *Scanner) regex() error {
	for s.peek() != '/' && !s.isAtEnd() {
		if s.peek() == '\n' {
			return errors.New("Unterminated RegEx.\n")
		}
		s.advance()
	}
	if s.isAtEnd() {
		return errors.New("Unterminated RegEx.\n")
	}
	s.advance() // Closing /

	s.addToken(REGEX)
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
