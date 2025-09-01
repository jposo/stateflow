package stateflow

import "github.com/jposo/stateflow/tools"

var alphabet string = "abcdefghijklmnoprstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var numeric string = "0123456789"

func BuildFsm() *tools.Fsm {
	fsm := tools.Fsm{}
	start := tools.State{}
	fsm.CurrentState = &start

	start.AddNewTransition('(', &start)
	start.AddNewTransition(')', &start)
	start.AddNewTransition('{', &start)
	start.AddNewTransition('}', &start)
	// ->
	dash := start.AddNewTransition('-', nil)
	dash.AddNewTransition('>', &start)
	// <-
	less := start.AddNewTransition('<', nil)
	less.AddNewTransition('-', &start)
	// Comments
	slash := start.AddNewTransition('/', nil)
	slash.AddNewTransition('/', &start)
	// Strings
	start.AddNewTransition('"', nil)

	start.AddNewTransition(';', &start)
	start.AddNewTransition('\n', &start)

	for _, whitepsace := range []rune{' ', '\t', '\r'} {
		start.AddNewTransition(whitepsace, &start)
	}

	var alphanumericStates []*tools.State
	for _, alphanumeric := range alphabet + numeric {
		state := tools.State{}
		state.AddNewTransition(alphanumeric, &state)
		alphanumericStates = append(alphanumericStates, &state)
	}

	for _, alpha := range alphabet {
		alphaState := start.AddNewTransition(alpha, &start)
		for _, alpanumeric := range alphabet + numeric {
			alphaState.AddNewTransition(alpanumeric, alphaState)
		}
	}

	return &fsm
}
