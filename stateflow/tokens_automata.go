package stateflow

type State int

const (
	StartState State = iota

	SlashState
	DashState
	LessState
	AlphabetState
	AlphanumericState
	OpenParenState
	CloseParenState
	OpenBraceState
	CloseBraceState

	InCommentState
	InIdentifierState
	InStringState

	ValidState
)

var transitions = map[State]map[byte]State{
	StartState: {
		'-': DashState,
		'<': LessState,
		'"': ValidState,
		'(': OpenParenState,
		')': CloseParenState,
		'{': OpenBraceState,
		'}': CloseBraceState,
	},
	DashState: {
		'>': ValidState,
	},
	LessState: {
		'-': ValidState,
	},
	SlashState: {
		'/': InCommentState,
		// Regex?
	},
}
