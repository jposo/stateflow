package stateflow

import (
	"testing"
)

func TestScannerBasicTokens(t *testing.T) {
	source := "dfa state initial final on when or fn str"
	scanner := Scanner{Source: []byte(source)}
	tokens, errors := scanner.ScanTokens()

	if len(errors) != 0 {
		t.Fatalf("Expected no errors, got %d", len(errors))
	}

	expectedTokens := []TokenType{BOF, DFA, STATE, INITIAL, FINAL, ON, WHEN, OR, FUNCTION, STRING, EOF}
	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, expected := range expectedTokens {
		if tokens[i].tokenType != expected {
			t.Errorf("Token %d: expected %v, got %v", i, expected, tokens[i].tokenType)
		}
	}
}

func TestScannerArrows(t *testing.T) {
	source := "-> <-"
	scanner := Scanner{Source: []byte(source)}
	tokens, errors := scanner.ScanTokens()

	if len(errors) != 0 {
		t.Fatalf("Expected no errors, got %d", len(errors))
	}

	expectedTokens := []TokenType{BOF, ARROW_RIGHT, ARROW_LEFT, EOF}
	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, expected := range expectedTokens {
		if tokens[i].tokenType != expected {
			t.Errorf("Token %d: expected %v, got %v", i, expected, tokens[i].tokenType)
		}
	}
}

func TestScannerParenthesesAndBraces(t *testing.T) {
	source := "(){}(),;"
	scanner := Scanner{Source: []byte(source)}
	tokens, errors := scanner.ScanTokens()

	if len(errors) != 0 {
		t.Fatalf("Expected no errors, got %d", len(errors))
	}

	expectedTokens := []TokenType{BOF, LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, LEFT_PAREN, RIGHT_PAREN, COMMA, SEMICOLON, EOF}
	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, expected := range expectedTokens {
		if tokens[i].tokenType != expected {
			t.Errorf("Token %d: expected %v, got %v", i, expected, tokens[i].tokenType)
		}
	}
}

func TestScannerComments(t *testing.T) {
	source := "dfa // this is a comment\nstate"
	scanner := Scanner{Source: []byte(source)}
	tokens, errors := scanner.ScanTokens()

	if len(errors) != 0 {
		t.Fatalf("Expected no errors, got %d", len(errors))
	}

	// Comments should be skipped, so we expect: BOF, DFA, STATE, EOF
	expectedTokens := []TokenType{BOF, DFA, STATE, EOF}
	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, expected := range expectedTokens {
		if tokens[i].tokenType != expected {
			t.Errorf("Token %d: expected %v, got %v", i, expected, tokens[i].tokenType)
		}
	}
}

func TestScannerLineTracking(t *testing.T) {
	source := "dfa\nstate\nfinal"
	scanner := Scanner{Source: []byte(source)}
	tokens, errors := scanner.ScanTokens()

	if len(errors) != 0 {
		t.Fatalf("Expected no errors, got %d", len(errors))
	}

	// Check that line numbers are correctly tracked
	if tokens[1].line != 1 {
		t.Errorf("Token 1 (dfa): expected line 1, got %d", tokens[1].line)
	}
	if tokens[2].line != 2 {
		t.Errorf("Token 2 (state): expected line 2, got %d", tokens[2].line)
	}
	if tokens[3].line != 3 {
		t.Errorf("Token 3 (final): expected line 3, got %d", tokens[3].line)
	}
}
