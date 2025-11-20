package stateflow

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitStateDeclStmt(stmt StateDecl) (any, error)
	VisitTransDeclStmt(stmt TransDecl) (any, error)
}

type StateDecl struct {
	stateType Token
	name      Token
}

func (s StateDecl) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitStateDeclStmt(s)
}

type TransDecl struct {
	symbol     Token
	fromState  Token
	toState    Token
	conditions []Condition
}

func (t TransDecl) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitTransDeclStmt(t)
}
