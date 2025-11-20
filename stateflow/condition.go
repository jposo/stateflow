package stateflow

type Condition interface {
	Accept(visitor ConditionVisitor) (any, error)
}

type ConditionVisitor interface {
	VisitStringConditionCondition(condition StringCondition) (any, error)
	VisitRegexConditionCondition(condition RegexCondition) (any, error)
}

type StringCondition struct {
	value string
}

func (s StringCondition) Accept(visitor ConditionVisitor) (any, error) {
	return visitor.VisitStringConditionCondition(s)
}

type RegexCondition struct {
	pattern string
}

func (r RegexCondition) Accept(visitor ConditionVisitor) (any, error) {
	return visitor.VisitRegexConditionCondition(r)
}
