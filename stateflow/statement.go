package stateflow

type Statement interface {
	Accept(visitor StatementVisitor) (any, error)
}

type StatementVisitor interface {
	VisitAssignmentStatement(statement Assignment) (any, error)
}

type Assignment struct {
	target Token
	source Token
}

func (a Assignment) Accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitAssignmentStatement(a)
}
