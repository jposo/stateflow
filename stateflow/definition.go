package stateflow

type Definition interface {
	Accept(visitor DefinitionVisitor) (any, error)
}

type DefinitionVisitor interface {
	VisitAutomatonDefDefinition(definition AutomatonDef) (any, error)
	VisitFunctionDefDefinition(definition FunctionDef) (any, error)
}

type AutomatonDef struct {
	autType Token
	name    Token
	stmts   []Stmt
}

func (a AutomatonDef) Accept(visitor DefinitionVisitor) (any, error) {
	return visitor.VisitAutomatonDefDefinition(a)
}

type FunctionDef struct {
	name       Token
	params     []Token
	statements []Statement
}

func (f FunctionDef) Accept(visitor DefinitionVisitor) (any, error) {
	return visitor.VisitFunctionDefDefinition(f)
}
