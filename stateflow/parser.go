package stateflow

type Parser struct {
	tokens   []Token
	current  int
	hadError bool
}
