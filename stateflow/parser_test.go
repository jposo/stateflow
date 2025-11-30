package stateflow

import (
	"testing"
)

// Helper function to create a scanner and get tokens
func getTokens(source string) []Token {
	scanner := Scanner{Source: []byte(source)}
	tokens, _ := scanner.ScanTokens()
	return tokens
}

// Test 1: Valid DFA with single final state
func TestParseSimpleDFA(t *testing.T) {
	source := `dfa simple {
		final q0;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) != 1 {
		t.Errorf("Expected 1 definition, got %d", len(defs))
	}

	if autoDef, ok := defs[0].(*AutomatonDef); ok {
		if autoDef.name.lexeme != "simple" {
			t.Errorf("Expected automaton name 'simple', got '%s'", autoDef.name.lexeme)
		}
		if autoDef.autType.tokenType != DFA {
			t.Errorf("Expected DFA type, got %v", autoDef.autType.tokenType)
		}
	} else {
		t.Error("Expected AutomatonDef")
	}
}

// Test 2: Valid NFA with states
// func TestParseNFA(t *testing.T) {
// 	source := `nfa machine {
// 		initial s0;
// 		state s1;
// 		final s2;

// 		on s0 -> s1 when "a";
// 		on s1 -> s2 when "b";
// 	}`
// 	tokens := getTokens(source)
// 	parser := Parser{Tokens: tokens}

// 	defs, err := parser.Parse()

// 	if err != nil {
// 		t.Errorf("Expected no error, got: %v", err)
// 	}

// 	if len(defs) != 1 {
// 		t.Errorf("Expected 1 definition, got %d", len(defs))
// 	}

// 	if autoDef, ok := defs[0].(*AutomatonDef); ok {
// 		if autoDef.autType.tokenType != NFA {
// 			t.Errorf("Expected NFA type, got %v", autoDef.autType.tokenType)
// 		}
// 		if len(autoDef.stmts) != 5 {
// 			t.Errorf("Expected 5 statements, got %d", len(autoDef.stmts))
// 		}
// 	}
// }

// Test 3: Automaton with transitions
func TestParseTransitions(t *testing.T) {
	source := `dfa counter {
		initial q0;
		state q1;
		final q2;
		
		on q0 -> q1 when "inc";
		on q1 -> q2 when "done";
		on q2 -> q2 when "reset";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) == 0 {
		t.Fatal("Expected at least 1 definition")
	}

	if autoDef, ok := defs[0].(*AutomatonDef); ok {
		transitionCount := 0
		for _, stmt := range autoDef.stmts {
			if _, ok := stmt.(*TransDecl); ok {
				transitionCount++
			}
		}
		if transitionCount != 3 {
			t.Errorf("Expected 3 transitions, got %d", transitionCount)
		}
	}
}

// Test 4: Function definition
func TestParseFunction(t *testing.T) {
	source := `fn process(input) {
		output <- input;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) != 1 {
		t.Errorf("Expected 1 definition, got %d", len(defs))
	}

	if funcDef, ok := defs[0].(*FunctionDef); ok {
		if funcDef.name.lexeme != "process" {
			t.Errorf("Expected function name 'process', got '%s'", funcDef.name.lexeme)
		}
		if len(funcDef.params) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(funcDef.params))
		}
		if funcDef.params[0].lexeme != "input" {
			t.Errorf("Expected parameter 'input', got '%s'", funcDef.params[0].lexeme)
		}
	}
}

// Test 5: Function with multiple parameters
func TestParseFunctionMultipleParams(t *testing.T) {
	source := `dfa calc {
		initial start;

		on start -> start when "1";
	}

	fn add(a, b, c) {
		calc <- a;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	if funcDef, ok := defs[0].(*FunctionDef); ok {
		if len(funcDef.params) != 3 {
			t.Errorf("Expected 3 parameters, got %d", len(funcDef.params))
		}
	}
}

// Test 6: Multiple definitions
func TestParseMultipleDefinitions(t *testing.T) {
	source := `dfa first {
		initial q0;
		on q0 -> q0 when "tick";
	}
	
	dfa second {
		initial s0;
		on s0 -> s0 when "tock";
	}
	
	fn myFunc(x) {
		first <- x;
		second <- x;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) != 3 {
		t.Errorf("Expected 3 definitions, got %d", len(defs))
	}
}

// Test 7: Duplicate automaton name (should error)
func TestErrorDuplicateAutomatonName(t *testing.T) {
	source := `dfa machine {
		initial q0;
	}
	
	dfa machine {
		initial s0;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for duplicate automaton name")
	}
}

// Test 8: Empty automaton (should error)
func TestErrorEmptyAutomaton(t *testing.T) {
	source := `dfa empty {
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for empty automaton")
	}
}

// Test 9: Duplicate state names (should error)
func TestErrorDuplicateStateNames(t *testing.T) {
	source := `dfa test {
		initial q0;
		state q0;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for duplicate state names")
	}
}

// Test 10: Undefined state reference (should error)
func TestErrorUndefinedStateReference(t *testing.T) {
	source := `dfa test {
		initial q0;
		on q0 -> undefined when "a";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for undefined state reference")
	}
}

// Test 11: Invalid parameter usage (should error)
func TestErrorInvalidParameterUsage(t *testing.T) {
	source := `fn test(a) {
		b <- undefined;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for undefined parameter")
	}
}

// Test 12: Multiple transitions with OR condition
func TestParseTransitionWithOR(t *testing.T) {
	source := `dfa test {
		initial q0;
		final q1;
		
		on q0 -> q1 when "a" or "b" or /[0-9]/;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) > 0 {
		if autoDef, ok := defs[0].(*AutomatonDef); ok {
			if len(autoDef.stmts) != 3 {
				t.Errorf("Expected 3 statements (initial + final + transition), got %d", len(autoDef.stmts))
			}
		}
	}
}

// Test 13: Duplicate function name (should error)
func TestErrorDuplicateFunctionName(t *testing.T) {
	source := `fn process(x) {
		y <- x;
	}
	
	fn process(a) {
		b <- a;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for duplicate function name")
	}
}

// Test 14: Final state with outgoing transition (should error)
func TestErrorFinalStateWithTransition(t *testing.T) {
	source := `dfa test {
		initial q0;
		final q1;
		
		on q1 -> q0 when "a";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for final state with outgoing transition")
	}
}

// Test 15: DFA with non-deterministic transitions (should error)
func TestErrorDFANonDeterministic(t *testing.T) {
	source := `dfa test {
		initial q0;
		state q1;
		final q2;
		
		on q0 -> q1 when "a";
		on q0 -> q2 when "a";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for non-deterministic DFA")
	}
}

// Test 16: Symbol table registration
func TestSymbolTableRegistration(t *testing.T) {
	source := `dfa myDFA {
		initial q0;
		on q0 -> q0 when "loop";
	}
	
	fn main(param) {
		myDfa <- param;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check automaton registration
	automaton := parser.SymbolTable.Lookup("myDFA")
	if automaton == nil {
		t.Error("Expected automaton 'myDFA' in symbol table")
	} else if automaton.Type != SymbolAutomaton {
		t.Errorf("Expected SymbolAutomaton type, got %s", automaton.Type)
	}

	// Check function registration
	function := parser.SymbolTable.Lookup("main")
	if function == nil {
		t.Error("Expected function 'main' in symbol table")
	} else if function.Type != SymbolFunction {
		t.Errorf("Expected SymbolFunction type, got %s", function.Type)
	}
}

// Test 17: Symbol table metadata
func TestSymbolTableMetadata(t *testing.T) {
	source := `dfa myDFA {
		initial q0;
		on q0 -> q0 when "loop";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	automaton := parser.SymbolTable.Lookup("myDFA")
	if automaton == nil {
		t.Fatal("Expected automaton in symbol table")
	}

	autType, ok := automaton.Metadata["automatonType"]
	if !ok {
		t.Error("Expected 'automatonType' in metadata")
	}

	if autType != DFA {
		t.Errorf("Expected DFA in metadata, got %v", autType)
	}
}

// Test 18: Valid complex automaton
func TestParseComplexAutomaton(t *testing.T) {
	source := `dfa complexDFA {
		initial start;
		state middle;
		state error;
		final success;
		
		on start -> middle when "go";
		on middle -> success when "done";
		on middle -> error when "fail";
		on error -> middle when "retry";
		on success -> success when "reset";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) == 0 {
		t.Fatal("Expected 1 definition")
	}

	if autoDef, ok := defs[0].(*AutomatonDef); ok {
		// Should have 4 state declarations + 5 transitions = 9 statements
		expectedStmts := 9
		if len(autoDef.stmts) != expectedStmts {
			t.Errorf("Expected %d statements, got %d", expectedStmts, len(autoDef.stmts))
		}
	}
}

// Test 19: Multiple initial states (should error)
func TestErrorMultipleInitialStates(t *testing.T) {
	source := `dfa test {
		initial q0;
		initial q1;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	if err == nil {
		t.Error("Expected error for multiple initial states")
	}
}

// Test 20: DFA state without transitions (non-final, should error)
func TestErrorDFAStateWithoutTransitions(t *testing.T) {
	source := `dfa test {
		initial q0;
		state q1;
		final q2;
		
		on q0 -> q1 when "a";
		on q1 -> q2 when "b";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	// q1 is non-final and has a transition, so this should PASS
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Test 20b: DFA state without transitions (non-final, should error)
func TestErrorDFAStateWithoutOutgoingTransitions(t *testing.T) {
	source := `dfa test {
		initial q0;
		state q1;
		final q2;
		
		on q0 -> q1 when "a";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	_, err := parser.Parse()

	// q1 is non-final but has no outgoing transitions - this should ERROR
	if err == nil {
		t.Error("Expected error for non-final state without outgoing transitions")
	}
}

// Test 21: NFA allows non-determinism
// func TestNFANonDeterminism(t *testing.T) {
// 	source := `nfa test {
// 		initial q0;
// 		state q1;
// 		final q2;

// 		on q0 -> q1 when "a";
// 		on q0 -> q2 when "a";
// 	}`
// 	tokens := getTokens(source)
// 	parser := Parser{Tokens: tokens}

// 	defs, err := parser.Parse()

// 	if err != nil {
// 		t.Errorf("Expected no error for NFA with non-determinism, got: %v", err)
// 	}

// 	if len(defs) == 0 {
// 		t.Fatal("Expected definitions")
// 	}
// }

// Test 22: Transition with regex condition
func TestParseRegexCondition(t *testing.T) {
	source := `dfa test {
		initial q0;
		final q1;
		
		on q0 -> q1 when /[a-z]+/;
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(defs) > 0 {
		if autoDef, ok := defs[0].(*AutomatonDef); ok {
			if len(autoDef.stmts) != 3 { // initial + final + transition
				t.Errorf("Expected 3 statements, got %d", len(autoDef.stmts))
			}
		}
	}
}

// Test 23: Function with multiple statements
// func TestParseFunctionMultipleStatements(t *testing.T) {
// 	source := `fn process(a, b, c) {
// 		x <- a;
// 		y <- b;
// 		z <- c;
// 	}`
// 	tokens := getTokens(source)
// 	parser := Parser{Tokens: tokens}

// 	defs, err := parser.Parse()

// 	if err != nil {
// 		t.Errorf("Expected no error, got: %v", err)
// 	}

// 	if funcDef, ok := defs[0].(*FunctionDef); ok {
// 		if len(funcDef.statements) != 3 {
// 			t.Errorf("Expected 3 statements, got %d", len(funcDef.statements))
// 		}
// 	}
// }

// Test 24: Empty parameter function
func TestParseFunctionNoParams(t *testing.T) {
	source := `fn noParams() {
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if funcDef, ok := defs[0].(*FunctionDef); ok {
		if len(funcDef.params) != 0 {
			t.Errorf("Expected 0 parameters, got %d", len(funcDef.params))
		}
	}
}

// Test 25: Automaton with all state types
func TestParseAllStateTypes(t *testing.T) {
	source := `dfa test {
		initial q0;
		state q1;
		state q2;
		final q3;

		on q0 -> q1 when "go";
		on q1 -> q2 when "continue";
		on q2 -> q3 when "finish";
	}`
	tokens := getTokens(source)
	parser := Parser{Tokens: tokens}

	defs, err := parser.Parse()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if autoDef, ok := defs[0].(*AutomatonDef); ok {
		stateCount := 0
		for _, stmt := range autoDef.stmts {
			if _, ok := stmt.(*StateDecl); ok {
				stateCount++
			}
		}
		if stateCount != 4 {
			t.Errorf("Expected 4 state declarations, got %d", stateCount)
		}
	}
}
