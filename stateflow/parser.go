package stateflow

import (
	"slices"
)

// Symbol represents an entry in the symbol table
type Symbol struct {
	Name     string
	Type     SymbolType
	Token    *Token
	Metadata map[string]any // For storing additional info (params, states, etc.)
}

type SymbolType string

const (
	SymbolAutomaton SymbolType = "Automaton"
	SymbolFunction  SymbolType = "Function"
	SymbolState     SymbolType = "State"
	SymbolParam     SymbolType = "Parameter"
)

// SymbolTable manages symbol declarations and scoping
type SymbolTable struct {
	scopes []map[string]*Symbol // Stack of scopes
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		scopes: []map[string]*Symbol{make(map[string]*Symbol)}, // Global scope
	}
}

func (st *SymbolTable) Define(name string, symbol *Symbol) error {
	current := st.scopes[len(st.scopes)-1]
	if _, exists := current[name]; exists {
		return ParseError{symbol.Token, "Symbol '" + name + "' already defined in this scope."}
	}
	current[name] = symbol
	return nil
}

func (st *SymbolTable) Lookup(name string) *Symbol {
	// Search from innermost to outermost scope
	for i := len(st.scopes) - 1; i >= 0; i-- {
		if sym, exists := st.scopes[i][name]; exists {
			return sym
		}
	}
	return nil
}

func (st *SymbolTable) PushScope() {
	st.scopes = append(st.scopes, make(map[string]*Symbol))
}

func (st *SymbolTable) PopScope() {
	if len(st.scopes) > 1 {
		st.scopes = st.scopes[:len(st.scopes)-1]
	}
}

type Parser struct {
	Tokens      []Token
	current     int
	SymbolTable *SymbolTable
}

func (p *Parser) Parse() ([]Definition, error) {
	if p.SymbolTable == nil {
		p.SymbolTable = NewSymbolTable()
	}

	if !p.match(BOF) {
		return nil, ParseError{p.peek(), "Expect BOF at start of program."}
	}
	var definitions []Definition
	for !p.isAtEnd() && !p.check(EOF) {
		definition, err := p.definition()
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}

	if _, err := p.consume(EOF, "Expect EOF at end of program."); err != nil {
		return nil, err
	}
	return definitions, nil
}

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, func(t TokenType) bool { return p.check(t) }) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) peek() *Token {
	return &p.Tokens[p.current]
}

func (p *Parser) previous() *Token {
	return &p.Tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.Tokens)
}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return nil, ParseError{p.peek(), message}
}

func (p *Parser) definition() (Definition, error) {
	if p.check(DFA) || p.check(NFA) {
		return p.automatonDef()
	}
	if p.check(FUNCTION) {
		return p.functionDef()
	}
	return nil, ParseError{p.peek(), "Expect automaton or function definition."}
}

func (p *Parser) automatonDef() (*AutomatonDef, error) {
	automatonType, err := p.automatonType()
	if err != nil {
		return nil, err
	}

	name, err := p.consume(IDENTIFIER, "Expect automaton name.")
	if err != nil {
		return nil, err
	}

	// Register automaton in symbol table
	if err := p.SymbolTable.Define(name.lexeme, &Symbol{
		Name:  name.lexeme,
		Type:  SymbolAutomaton,
		Token: name,
		Metadata: map[string]any{
			"automatonType": automatonType.tokenType,
		},
	}); err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before automaton body.")
	if err != nil {
		return nil, err
	}

	stmts, err := p.stmtList()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after automaton body.")
	if err != nil {
		return nil, err
	}

	// Validate that final states don't have outgoing transitions
	if err := p.validateAutomaton(stmts, automatonType.tokenType); err != nil {
		return nil, err
	}

	return &AutomatonDef{
		autType: *automatonType,
		name:    *name,
		stmts:   stmts,
	}, nil
}

// validateAutomaton checks various constraints on the automaton
func (p *Parser) validateAutomaton(stmts []Stmt, automatonType TokenType) error {
	// Check that at least one state is declared
	if err := p.validateNonEmptyAutomaton(stmts); err != nil {
		return err
	}

	// Check for duplicate state names
	if err := p.validateUniqueStates(stmts); err != nil {
		return err
	}

	// Check for duplicate initial states
	if err := p.validateUniqueInitialState(stmts); err != nil {
		return err
	}

	// Check that final states don't have outgoing transitions
	if err := p.validateFinalStates(stmts); err != nil {
		return err
	}

	// Check that all referenced states exist
	if err := p.validateStateReferences(stmts); err != nil {
		return err
	}

	// For DFAs, validate deterministic transition rules
	if automatonType == DFA {
		if err := p.validateDFATransitions(stmts); err != nil {
			return err
		}
	}

	return nil
}

// Checks DFA-specific transition constraints
func (p *Parser) validateDFATransitions(stmts []Stmt) error {
	// Collect all states (non-final states that need outgoing transitions)
	states := make(map[string]bool)
	finalStates := make(map[string]bool)

	for _, stmt := range stmts {
		if stateDecl, ok := stmt.(*StateDecl); ok {
			states[stateDecl.name.lexeme] = true
			if stateDecl.stateType.tokenType == FINAL {
				finalStates[stateDecl.name.lexeme] = true
			}
		}
	}

	// Track transitions: map[fromState]map[symbol]toState
	transitions := make(map[string]map[string]Token)

	for _, stmt := range stmts {
		if transDecl, ok := stmt.(*TransDecl); ok {
			fromState := transDecl.fromState.lexeme

			// Initialize map for this state if needed
			if transitions[fromState] == nil {
				transitions[fromState] = make(map[string]Token)
			}

			// Check each condition (symbol) in this transition
			for _, condition := range transDecl.conditions {
				var symbol string

				switch cond := condition.(type) {
				case StringCondition:
					symbol = cond.value
					if symbol == "\"\"" {
						return ParseError{
							&transDecl.fromState,
							"Empty string condition not allowed in state '" + fromState + "'.",
						}
					}
				case RegexCondition:
					symbol = cond.pattern
				}

				// Check for duplicate transition on same symbol from same state
				if _, exists := transitions[fromState][symbol]; exists {
					return ParseError{
						&transDecl.fromState,
						"Duplicate transition from state '" + fromState +
							"' on symbol '" + symbol + "'. " +
							"DFA cannot have multiple transitions for the same symbol from the same state.",
					}
				}

				// Record this transition
				transitions[fromState][symbol] = transDecl.toState
			}
		}
	}

	// Check that each non-final state has exactly one transition
	// (assuming the automaton has a finite alphabet that should be covered)
	for state := range states {
		if !finalStates[state] {
			transCount := len(transitions[state])
			if transCount == 0 {
				// Find the token for this state to provide better error context
				var stateToken *Token
				for _, stmt := range stmts {
					if stateDecl, ok := stmt.(*StateDecl); ok {
						if stateDecl.name.lexeme == state {
							stateToken = &stateDecl.name
							break
						}
					}
				}

				return ParseError{
					stateToken,
					"State '" + state + "' has no outgoing transitions. " +
						"DFA requires every non-final state to have transitions.",
				}
			}
		}
	}

	return nil
}

// Checks that final states don't have outgoing transitions
func (p *Parser) validateFinalStates(stmts []Stmt) error {
	// Collect all final states
	finalStates := make(map[string]bool)
	for _, stmt := range stmts {
		if stateDecl, ok := stmt.(*StateDecl); ok {
			if stateDecl.stateType.tokenType == FINAL {
				finalStates[stateDecl.name.lexeme] = true
			}
		}
	}

	// Check if any final state has an outgoing transition
	// that isn't itself
	for _, stmt := range stmts {
		if transDecl, ok := stmt.(*TransDecl); ok {
			if finalStates[transDecl.fromState.lexeme] &&
				transDecl.toState.lexeme != transDecl.fromState.lexeme {
				return ParseError{
					&transDecl.fromState,
					"Final state '" + transDecl.fromState.lexeme + "' cannot have outgoing transitions.",
				}
			}
		}
	}

	return nil
}

// Checks that there is only one initial state
func (p *Parser) validateUniqueInitialState(stmts []Stmt) error {
	var initialState *StateDecl

	for _, stmt := range stmts {
		if stateDecl, ok := stmt.(*StateDecl); ok {
			if stateDecl.stateType.tokenType == INITIAL {
				if initialState != nil {
					return ParseError{
						&stateDecl.name,
						"Duplicate initial state '" + stateDecl.name.lexeme + "'. " +
							"Automaton already has initial state '" + initialState.name.lexeme + "'.",
					}
				}
				initialState = stateDecl
			}
		}
	}

	return nil
}

// Checks that automaton is not empty (has at least one state)
func (p *Parser) validateNonEmptyAutomaton(stmts []Stmt) error {
	hasState := false
	for _, stmt := range stmts {
		if _, ok := stmt.(*StateDecl); ok {
			hasState = true
			break
		}
	}
	if !hasState {
		return ParseError{
			p.peek(),
			"Automaton must declare at least one state.",
		}
	}
	return nil
}

// Checks that all state names within an automaton are unique
func (p *Parser) validateUniqueStates(stmts []Stmt) error {
	states := make(map[string]*Token)
	for _, stmt := range stmts {
		if stateDecl, ok := stmt.(*StateDecl); ok {
			if existing, found := states[stateDecl.name.lexeme]; found {
				return ParseError{
					&stateDecl.name,
					"Duplicate state declaration '" + stateDecl.name.lexeme + "'. " +
						"State was already declared at line " + string(rune(existing.line)),
				}
			}
			states[stateDecl.name.lexeme] = &stateDecl.name
		}
	}
	return nil
}

// Checks that all states referenced in transitions exist
func (p *Parser) validateStateReferences(stmts []Stmt) error {
	// First collect all declared states
	declaredStates := make(map[string]bool)
	for _, stmt := range stmts {
		if stateDecl, ok := stmt.(*StateDecl); ok {
			declaredStates[stateDecl.name.lexeme] = true
		}
	}

	// Check all transitions reference valid states
	for _, stmt := range stmts {
		if transDecl, ok := stmt.(*TransDecl); ok {
			if !declaredStates[transDecl.fromState.lexeme] {
				return ParseError{
					&transDecl.fromState,
					"Transition references undefined state '" + transDecl.fromState.lexeme + "'.",
				}
			}
			if !declaredStates[transDecl.toState.lexeme] {
				return ParseError{
					&transDecl.toState,
					"Transition references undefined state '" + transDecl.toState.lexeme + "'.",
				}
			}
		}
	}

	return nil
}

func (p *Parser) automatonType() (*Token, error) {
	if p.match(DFA) {
		return p.previous(), nil
	}
	if p.match(NFA) {
		return p.previous(), nil
	}
	return nil, ParseError{p.peek(), "Expect 'dfa' or 'nfa'."}
}

func (p *Parser) stmtList() ([]Stmt, error) {
	var stmts []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.stmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)

		_, err = p.consume(SEMICOLON, "Expect ';' after statement.")
		if err != nil {
			return nil, err
		}
	}
	return stmts, nil
}

// Statements for automatas
func (p *Parser) stmt() (Stmt, error) {
	if p.check(INITIAL) || p.check(STATE) || p.check(FINAL) {
		return p.stateDecl()
	}
	if p.check(ON) {
		return p.transDecl()
	}
	return nil, ParseError{p.peek(), "Expect state declaration or transition declaration."}
}

func (p *Parser) stateDecl() (*StateDecl, error) {
	declType, err := p.stateDeclType()
	if err != nil {
		return nil, err
	}
	name, err := p.consume(IDENTIFIER, "Expect state name.")
	if err != nil {
		return nil, err
	}

	return &StateDecl{
		stateType: *declType,
		name:      *name,
	}, nil
}

func (p *Parser) stateDeclType() (*Token, error) {
	if p.match(INITIAL) {
		return p.previous(), nil
	}
	if p.match(STATE) {
		return p.previous(), nil
	}
	if p.match(FINAL) {
		return p.previous(), nil
	}
	return nil, ParseError{p.peek(), "Expect 'initial', 'state', or 'final'."}
}

func (p *Parser) transDecl() (*TransDecl, error) {
	_, err := p.consume(ON, "Expect 'on'.")
	if err != nil {
		return nil, err
	}

	fromState, err := p.consume(IDENTIFIER, "Expect source state name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(ARROW_RIGHT, "Expect '->' in transition.")
	if err != nil {
		return nil, err
	}

	toState, err := p.consume(IDENTIFIER, "Expect target state name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(WHEN, "Expect 'when' before conditions.")
	if err != nil {
		return nil, err
	}

	conditions, err := p.conditionList()
	if err != nil {
		return nil, err
	}

	return &TransDecl{
		fromState:  *fromState,
		toState:    *toState,
		conditions: conditions,
	}, nil
}

func (p *Parser) conditionList() ([]Condition, error) {
	var conditions []Condition

	condition, err := p.condition()
	if err != nil {
		return nil, err
	}
	conditions = append(conditions, condition)

	for p.match(OR) {
		condition, err := p.condition()
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}

func (p *Parser) condition() (Condition, error) {
	if p.match(STRING_LITERAL) {
		return StringCondition{value: p.previous().lexeme}, nil
	}
	if p.match(REGEX) {
		return RegexCondition{pattern: p.previous().lexeme}, nil
	}
	return nil, ParseError{p.peek(), "Expect string or regex condition."}
}

func (p *Parser) functionDef() (*FunctionDef, error) {
	_, err := p.consume(FUNCTION, "Expect 'fn'.")
	if err != nil {
		return nil, err
	}

	name, err := p.consume(IDENTIFIER, "Expect function name.")
	if err != nil {
		return nil, err
	}

	// Register function in global symbol table
	if err := p.SymbolTable.Define(name.lexeme, &Symbol{
		Name:  name.lexeme,
		Type:  SymbolFunction,
		Token: name,
		Metadata: map[string]any{
			"params": []string{},
		},
	}); err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, "Expect '(' after function name.")
	if err != nil {
		return nil, err
	}

	params, err := p.paramsDef()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_BRACE, "Expect '{' before function body.")
	if err != nil {
		return nil, err
	}

	// Create new scope for function body
	p.SymbolTable.PushScope()

	// Register parameters in function scope
	paramNames := []string{}
	for _, param := range params {
		if err := p.SymbolTable.Define(param.lexeme, &Symbol{
			Name:  param.lexeme,
			Type:  SymbolParam,
			Token: &param,
		}); err != nil {
			p.SymbolTable.PopScope()
			return nil, err
		}
		paramNames = append(paramNames, param.lexeme)
	}

	// Update function metadata with parameter names
	funcSym := p.SymbolTable.Lookup(name.lexeme)
	if funcSym != nil {
		funcSym.Metadata["params"] = paramNames
	}

	statements, err := p.statementList()
	if err != nil {
		p.SymbolTable.PopScope()
		return nil, err
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after function body.")
	if err != nil {
		p.SymbolTable.PopScope()
		return nil, err
	}

	// Pop function scope
	p.SymbolTable.PopScope()

	return &FunctionDef{
		name:       *name,
		params:     params,
		statements: statements,
	}, nil

}

func (p *Parser) paramsDef() ([]Token, error) {
	if p.check(RIGHT_PAREN) {
		return []Token{}, nil
	}
	return p.paramList()
}

func (p *Parser) paramList() ([]Token, error) {
	var params []Token

	param, err := p.param()
	if err != nil {
		return nil, err
	}
	params = append(params, param)

	for p.match(COMMA) {
		param, err := p.param()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}

	return params, nil
}

func (p *Parser) param() (Token, error) {
	token, err := p.consume(IDENTIFIER, "Expect parameter name.")
	if err != nil {
		return Token{}, err
	}
	return *token, nil
}

func (p *Parser) statementList() ([]Statement, error) {
	var statements []Statement

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)

		_, err = p.consume(SEMICOLON, "Expect ';' after statement.")
		if err != nil {
			return nil, err
		}
	}
	return statements, nil
}

// Statement inside functions
func (p *Parser) statement() (Statement, error) {
	target, err := p.consume(IDENTIFIER, "Expect identifier.")
	if err != nil {
		return Assignment{}, err
	}

	_, err = p.consume(ARROW_LEFT, "Expect '<-' in assignment.")
	if err != nil {
		return Assignment{}, err
	}

	source, err := p.consume(IDENTIFIER, "Expect identifier after '<-'.")
	if err != nil {
		return Assignment{}, err
	}

	// Validate that source identifier is a declared parameter or variable
	if p.SymbolTable.Lookup(source.lexeme) == nil {
		return Assignment{}, ParseError{
			source,
			"Undefined variable or parameter '" + source.lexeme + "'. " +
				"Variable must be a function parameter.",
		}
	}

	return Assignment{
		target: *target,
		source: *source,
	}, nil
}
