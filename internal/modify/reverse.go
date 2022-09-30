package modify

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
    "mipt_formal/internal/tools"
)

func Reverse(m *nfa.Machine) {
    singleStart(m)
    singleAcceptance(m)

    invDelta := make([]common.Transition, 0, len(m.Delta))
    for to := range m.Delta {
        for by, dst := range m.Delta[to] {
            for from := range dst {
                invDelta = append(invDelta, common.Transition{
                    From: from,
                    To:   common.State(to),
                    By:   by,
                })
            }
        }
    }

    rev := nfa.New(m.Final.AsSlice(), m.Start.AsSlice(), invDelta)

    m.Delta = rev.Delta
    m.Start = rev.Start
    m.Final = rev.Final
}

func singleStart(m *nfa.Machine) {
    if m.Start.Size() == 1 {
        return
    }

    newStart := m.AddState()
    for s := range m.Start {
        if m.Final.Has(s) {
            m.Final.Insert(newStart)
        }
        for by, dst := range m.Delta[s] {
            for to := range dst {
                m.AddTransition(newStart, to, by)
            }
        }
    }

    mask := removeStates(m, m.Start)
    if alias, found := mask[newStart]; found {
        newStart = alias
    }
    m.Start = tools.NewSet[common.State](newStart)
}

// singleAcceptance collapse all acceptances into one
// constraint: input machine has single start state
func singleAcceptance(m *nfa.Machine) {
    if m.Final.Size() == 1 {
        return
    }

    var newFinal common.State

    // First, check if single start is acceptance
    useStart := false
    for s := range m.Start {
        if m.Final.Has(s) {
            useStart = true
            newFinal = s
            m.Final.Delete(newFinal)
        }
    }

    if !useStart {
        newFinal = m.AddState()
    }

    for from := range m.Delta {
        for by, dst := range m.Delta[from] {
            for to := range dst {
                if m.Final.Has(to) {
                    if m.Start.Has(to) {
                        m.Start.Insert(newFinal)
                    }
                    m.AddTransition(common.State(from), newFinal, by)
                }
            }
        }
    }

    mask := removeStates(m, m.Final)
    if alias, found := mask[newFinal]; found {
        newFinal = alias
    }
    m.Final = tools.NewSet[common.State](newFinal)
}

func removeStates(m *nfa.Machine, toRemove tools.Set[common.State]) map[common.State]common.State {
    mask := calculateRemovingMask(m, toRemove)
    removeIngoingEdges(m, toRemove)
    applyMask(m, mask)
    clearFinals(m, toRemove, mask)
    return mask
}

func removeIngoingEdges(m *nfa.Machine, to tools.Set[common.State]) {
    for from := range m.Delta {
        for _, dst := range m.Delta[from] {
            for t := range dst {
                if to.Has(t) {
                    dst.Delete(t)
                }
            }
        }
    }
}
