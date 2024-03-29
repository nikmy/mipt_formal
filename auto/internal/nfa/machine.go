package nfa

import (
    "fmt"
    "mipt_formal/tools"
    "strings"

    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/doa"
)

type Transitions map[common.Word]tools.Set[common.State]

func NewMachine(start []common.State, final []common.State, edges []common.Transition) *Machine {
    nStates := 0
    for _, t := range edges {
        if int(t.From) >= nStates {
            nStates = int(t.From) + 1
        }
        if int(t.To) >= nStates {
            nStates = int(t.To) + 1
        }
    }

    for _, s := range start {
        if int(s) >= nStates {
            nStates = int(s) + 1
        }
    }

    for _, f := range final {
        if int(f) >= nStates {
            nStates = int(f) + 1
        }
    }

    d := make([]Transitions, nStates)
    for s := common.State(0); s < common.State(nStates); s++ {
        d[s] = make(Transitions)
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
        Start: tools.NewSet[common.State](start...),
        Final: tools.NewSet[common.State](final...),
    }
}

type Machine struct {
    Delta []Transitions
    Start tools.Set[common.State]
    Final tools.Set[common.State]
}

func (m *Machine) NStates() int {
    return len(m.Delta)
}

func (m *Machine) AddState() common.State {
    newState := common.State(m.NStates())
    m.Delta = append(m.Delta, make(Transitions))
    return newState
}

func (m *Machine) AddTransition(from common.State, to common.State, by common.Word) bool {
    if _, has := m.Delta[from][by]; !has {
        m.Delta[from][by] = tools.NewSet[common.State]()
    }
    return m.Delta[from][by].Insert(to)
}

func (m *Machine) Equal(other *Machine) bool {
    if m == nil || other == nil {
        return m == other
    }
    if m.NStates() != other.NStates() {
        return false
    }
    if m.Start.Size() != other.Start.Size() {
        return false
    }
    if m.Final.Size() != other.Final.Size() {
        return false
    }

    for s := range other.Start {
        if !m.Start.Has(s) {
            return false
        }
    }

    for f := range other.Final {
        if !m.Final.Has(f) {
            return false
        }
    }

    for i := range m.Delta {
        for w, set := range m.Delta[i] {
            if _, found := other.Delta[i][w]; !found {
                return false
            }
            for state := range set {
                if !other.Delta[i][w].Has(state) {
                    return false
                }
            }
        }
    }

    return true
}

func (m *Machine) MarkAsFinal(state common.State) {
    m.Final.Insert(state)
}

func (m *Machine) DOA() string {
    if m == nil {
        return ""
    }

    var b strings.Builder
    b.Grow(doa.MinimalLength)

    b.WriteString(doa.Version)

    start := m.Start.AsSlice()
    starts := fmt.Sprintf("%v", start[0])
    if m.Start.Size() > 1 {
        for _, s := range start[1:] {
            starts += fmt.Sprintf(doa.StateConj+"%v", s)
        }
    }
    b.WriteString(fmt.Sprintf(doa.StartFormat, starts))

    final := m.Final.AsSlice()
    finals := fmt.Sprintf("%v", final[0])
    if m.Final.Size() > 1 {
        for _, f := range final[1:] {
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

func (m *Machine) Go(state []common.State, w common.Word) []common.State {
    slice := make([]common.State, 0)
    for _, s := range state {
        if next, can := m.Delta[s][w]; can {
            slice = append(slice, next.AsSlice()...)
        }
    }
    return slice
}
