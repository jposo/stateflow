package stateflow

type Statement interface {
	Accept(visitor StatementVisitor) (any, error)
}

type StatementVisitor interface {
	VisitCallStatement(statement Call) (any, error)
}

type Call struct {
	target Token
	input  Token
}

func (c Call) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitCallStatement(c)
}
