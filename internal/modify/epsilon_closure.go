package modify

import (
    "mipt_formal/internal/tools"

    "mipt_formal/internal/nfa"
)

var _ Modifier = EliminateEpsilonMoves

func EliminateEpsilonMoves(stateMachine *nfa.Machine) {
    buildTransitiveEpsilonClosure(stateMachine)
    compressAcceptances(stateMachine)
    addTransitiveEdges(stateMachine)
    removeEpsilonMoves(stateMachine)
}

func buildTransitiveEpsilonClosure(m *nfa.Machine) {
    newEdges := true
    for newEdges {
        newEdges = false
        for i, t := range m.Delta {
            for j := range t[nfa.Epsilon] {
                for k := range m.Delta[j][nfa.Epsilon] {
                    if m.AddTransition(nfa.State(i), k, nfa.Epsilon) {
                        newEdges = true
                    }
                }
            }
        }
    }
}

func compressAcceptances(m *nfa.Machine) {
    f := tools.NewSet[nfa.State](m.Final...)
    newAccept := true
    for newAccept {
        newAccept = false
        for i, t := range m.Delta {
            s := nfa.State(i)
            for j := range t[nfa.Epsilon] {
                if f.Has(j) {
                    if f.Insert(s) {
                        m.MarkAsFinal(s)
                        newAccept = true
                    }
                }
            }
        }
    }
}

func addTransitiveEdges(m *nfa.Machine) {
    for i, t := range m.Delta {
        for j := range t[nfa.Epsilon] {
            for w, tt := range m.Delta[j] {
                for k := range tt {
                    m.AddTransition(nfa.State(i), k, w)
                }
            }
        }
    }
}

func removeEpsilonMoves(m *nfa.Machine) {
    for i := range m.Delta {
        delete(m.Delta[i], nfa.Epsilon)
    }
}
