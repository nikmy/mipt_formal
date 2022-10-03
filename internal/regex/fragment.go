package regex

import (
    "mipt_formal/internal/common"
)

func newIntrusiveState(label common.Word, next ...*intrusiveState) *intrusiveState {
    if len(next) == 0 {
        next = nil
    }
    return &intrusiveState{
        label: label,
        next:  next,
    }
}

type intrusiveState struct {
    label common.Word
    next  []*intrusiveState
}

func (s *intrusiveState) precede(others ...*intrusiveState) {
    s.next = append(s.next, others...)
}

type fragment struct {
    Start  *intrusiveState
    Accept *intrusiveState
}
