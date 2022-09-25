package fsm

import (
    "fmt"
    "mipt_formal/internal/doa"
    "mipt_formal/internal/tools"
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

func (m *finiteStateMachine) Equal(other *finiteStateMachine) bool {
    checkSet := tools.NewSet[State]()
    for _, s := range m.start {
        checkSet.Insert(s)
    }
    for _, s := range other.start {
        if !checkSet.Has(s) {
            return false
        }
    }

    checkSet = tools.NewSet[State]()
    for _, s := range m.final {
        checkSet.Insert(s)
    }
    for _, s := range other.final {
        if !checkSet.Has(s) {
            return false
        }
    }

    for s, t := range other.delta {
        check, ok := m.delta[s]
        if !ok {
            return false
        }
        for word, to := range t {
            checkTo, ok := check[word]
            if !ok || checkTo != to {
                return false
            }
        }
    }
    return true
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
