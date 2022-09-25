package fsm

import (
    "fmt"
    "mipt_formal/internal/doa"
    "strings"
)

func NewNFA(start []State, final []State, edges []Transition) *NFA {
    d := make([]transitions, 0, len(edges)) // TODO
    return &NFA{
        delta: d,
        start: start,
        final: final,
    }
}

type transitions map[Word][]State

type NFA struct {
    delta []transitions
    start []State
    final []State
}

func (m *NFA) Equal(other *NFA) bool {
    panic("not implemented")
}

func (m *NFA) DOA() string {
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
        for word, states := range t {
            for _, state := range states {
                if word == "" {
                    word = doa.Epsilon
                }
                b.WriteString(fmt.Sprintf(doa.EdgeFormat, word, state))
            }
        }
    }

    b.WriteString(doa.End)
    return b.String()
}

func (m *NFA) Go(s State, w Word) []State {
    panic("not implemented")
}
