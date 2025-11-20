package stateflow

type Parser struct {
	Tokens   []Token
	current  int
	hadError bool
}

func (p *Parser) Parse() ([]Definition, error) {
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
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
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

	return &AutomatonDef{
		autType: *automatonType,
		name:    *name,
		stmts:   stmts,
	}, nil
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
	if p.match(STRING) {
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

	statements, err := p.statementList()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after function body.")
	if err != nil {
		return nil, err
	}

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

	return Assignment{
		target: *target,
		source: *source,
	}, nil
}
