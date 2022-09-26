package regex

import "mipt_formal/internal/nfa"

func NewIntrusiveState(label nfa.Word, next ...*IntrusiveState) *IntrusiveState {
    if len(next) == 0 {
        next = nil
    }
    return &IntrusiveState{
        label: label,
        next:  next,
    }
}

type IntrusiveState struct {
    label nfa.Word
    next  []*IntrusiveState
}

func (s *IntrusiveState) precede(others ...*IntrusiveState) {
    s.next = append(s.next, others...)
}

type fragment struct {
    Start  *IntrusiveState
    Accept *IntrusiveState
}
