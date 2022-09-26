package nfa

import (
    "fmt"
    "mipt_formal/internal/doa"
    "mipt_formal/internal/tools"
    "strings"
)

type transitions map[Word]tools.Set[State]

func New(start []State, final []State, edges []Transition) *Machine {
    nStates := 0
    for _, t := range edges {
        if int(t.From) >= nStates {
            nStates = int(t.From) + 1
        }
        if int(t.To) >= nStates {
            nStates = int(t.To) + 1
        }
    }

    d := make([]transitions, nStates)
    for s := State(0); s < State(nStates); s++ {
        d[s] = nil
    }

    for _, t := range edges {
        if d[t.From] == nil {
            d[t.From] = make(transitions, 1)
        }
        if d[t.From][t.By] == nil {
            d[t.From][t.By] = make(tools.Set[State])
        }
        d[t.From][t.By].Insert(t.To)
    }

    return &Machine{
        Delta: d,
        Start: start,
        Final: final,
    }
}

type Machine struct {
    Delta []transitions
    Start []State
    Final []State
}

func (m *Machine) AddTransition(from State, to State, by Word) bool {
    if _, has := m.Delta[from][by]; !has {
        m.Delta[from][by] = tools.NewSet[State]()
    }
    return m.Delta[from][by].Insert(to)
}

func (m *Machine) MarkAsFinal(state State) {
    m.Final = append(m.Final, state)
}

func (m *Machine) DOA() string {
    var b strings.Builder
    b.Grow(doa.MinimalLength)

    b.WriteString(doa.Version)

    starts := fmt.Sprintf("%v", m.Start[0])
    if len(m.Start) > 1 {
        for _, s := range m.Start[1:] {
            starts += fmt.Sprintf(doa.StateConj+"%v", s)
        }
    }
    b.WriteString(fmt.Sprintf(doa.StartFormat, starts))

    finals := fmt.Sprintf("%v", m.Final[0])
    if len(m.Final) > 1 {
        for _, f := range m.Final[1:] {
            finals += fmt.Sprintf(doa.StateConj+"%v", f)
        }
    }
    b.WriteString(fmt.Sprintf(doa.AcceptanceFormat, finals))

    b.WriteString(doa.Begin)

    for s, t := range m.Delta {
        b.WriteString(fmt.Sprintf(doa.StateFormat, s))
        for word, states := range t {
            for state := range states {
                if word == Epsilon {
                    word = doa.Epsilon
                }
                b.WriteString(fmt.Sprintf(doa.EdgeFormat, word, state))
            }
        }
    }

    b.WriteString(doa.End)
    return b.String()
}

func (m *Machine) Go(s State, w Word) []State {
    panic("not implemented")
}
