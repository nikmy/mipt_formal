package fsm

import (
    "fmt"
    "mipt_formal/internal/doa"
    "strings"
)

func New() *finiteStateMachine {
    panic("not implemented")
}

type transitions map[Word]State

type finiteStateMachine struct {
    delta map[State]transitions
    start []State
    final []State
}

func (m *finiteStateMachine) DOA() string {
    var b strings.Builder
    b.Grow(doa.MinimalLength)

    b.WriteString(doa.Version)

    starts := fmt.Sprintf("%v", m.start[0])
    if len(m.start) > 1 {
        for _, s := range m.start[1:] {
            starts += fmt.Sprintf(doa.StateConj+"%v", s)
        }
    }
    b.WriteString(fmt.Sprintf(doa.StartFormat, starts))

    finals := fmt.Sprintf("%v", m.final[0])
    if len(m.final) > 1 {
        for _, f := range m.final[1:] {
            finals += fmt.Sprintf(doa.StateConj+"%v", f)
        }
    }
    b.WriteString(fmt.Sprintf(doa.AcceptanceFormat, finals))

    b.WriteString(doa.Begin)

    for s, t := range m.delta {
        b.WriteString(fmt.Sprintf(doa.StateFormat, s))
        for word, state := range t {
            if word == "" {
                word = doa.Epsilon
            }
            b.WriteString(fmt.Sprintf(doa.EdgeFormat, word, state))
        }
    }

    b.WriteString(doa.End)
    return b.String()
}

func (m *finiteStateMachine) Go(s State, w Word) []State {
    panic("not implemented")
}
