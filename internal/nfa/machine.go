package nfa

import (
    "fmt"
    "strings"

    "mipt_formal/internal/common"
    "mipt_formal/internal/doa"
    "mipt_formal/internal/tools"
)

type Transitions map[common.Word]tools.Set[common.State]

func New(start []common.State, final []common.State, edges []common.Transition) *Machine {
    nStates := 0
    for _, t := range edges {
        if int(t.From) >= nStates {
            nStates = int(t.From) + 1
        }
        if int(t.To) >= nStates {
            nStates = int(t.To) + 1
        }
    }

    d := make([]Transitions, nStates)
    for s := common.State(0); s < common.State(nStates); s++ {
        d[s] = nil
    }

    for _, t := range edges {
        if d[t.From] == nil {
            d[t.From] = make(Transitions, 1)
        }
        if d[t.From][t.By] == nil {
            d[t.From][t.By] = make(tools.Set[common.State])
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
    Delta []Transitions
    Start []common.State
    Final []common.State
}

func (m *Machine) NStates() int {
    return len(m.Delta)
}

func (m *Machine) AddTransition(from common.State, to common.State, by common.Word) bool {
    if _, has := m.Delta[from][by]; !has {
        m.Delta[from][by] = tools.NewSet[common.State]()
    }
    return m.Delta[from][by].Insert(to)
}

func (m *Machine) MarkAsFinal(state common.State) {
    m.Final = append(m.Final, state)
}

func (m *Machine) DOA() string {
    if m == nil {
        return ""
    }

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
                if word == common.Epsilon {
                    word = doa.Epsilon
                }
                b.WriteString(fmt.Sprintf(doa.EdgeFormat, word, state))
            }
        }
    }

    b.WriteString(doa.End)
    return b.String()
}

func (m *Machine) Go(s common.State, w common.Word) []common.State {
    panic("not implemented")
}
