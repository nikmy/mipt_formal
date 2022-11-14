package modify

import (
    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/nfa"
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
        for i := range m.Delta {
            for j := range m.Delta[i][common.Epsilon] {
                for k := range m.Delta[j][common.Epsilon] {
                    if m.AddTransition(common.State(i), k, common.Epsilon) {
                        newEdges = true
                    }
                }
            }
        }
    }
}

func compressAcceptances(m *nfa.Machine) {
    newAccept := true
    for newAccept {
        newAccept = false
        for i := range m.Delta {
            from := common.State(i)
            for to := range m.Delta[i][common.Epsilon] {
                if m.Final.Has(to) {
                    if m.Final.Insert(from) {
                        newAccept = true
                    }
                }
            }
        }
    }
}

func addTransitiveEdges(m *nfa.Machine) {
    for i := range m.Delta {
        for j := range m.Delta[i][common.Epsilon] {
            for word := range m.Delta[j] {
                for k := range m.Delta[j][word] {
                    m.AddTransition(common.State(i), k, word)
                }
            }
        }
    }
}

func removeEpsilonMoves(m *nfa.Machine) {
    for state := range m.Delta {
        delete(m.Delta[state], common.Epsilon)
    }
}
