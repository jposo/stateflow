package tools

import "fmt"

type State struct {
	transitions map[rune]*State
}

func (s *State) IsValidTransition(event rune) bool {
	_, ok := s.transitions[event]
	return ok
}

func (s *State) AddNewTransition(event rune, pointTo *State) *State {
	state := State{}
	s.transitions[event] = &state
	if pointTo != nil {
		s.transitions[rune(0)] = pointTo
	}
	return &state
}

func (s *State) AddManyNewTransition(events []rune) []*State {
	var states []*State
	for _, event := range events {
		state := State{}
		s.transitions[event] = &state
		states = append(states, &state)
	}
	return states
}

func (s *State) AddTransition(event rune, state *State) {
	s.transitions[event] = state
}

type Fsm struct {
	CurrentState *State
}

func (f *Fsm) Trigger(event rune) error {
	ok := f.CurrentState.IsValidTransition(event)
	if !ok {
		return fmt.Errorf("Invalid transition.")
	}
	f.CurrentState = f.CurrentState.transitions[event]
	return nil
}
