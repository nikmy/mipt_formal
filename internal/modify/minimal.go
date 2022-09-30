package modify

import (
    "mipt_formal/internal/nfa"
)

// Minimize build DFA with minimal number of states
// From Janusz A. (John) Brzozowski, 1962
func Minimize(m *nfa.Machine) {
    bSteps := []Modifier{
        EliminateEpsilonMoves, RemoveUnusedStates,
        Reverse, Determine, Reverse, Determine,
        RemoveUnusedStates,
    }
    for _, step := range bSteps {
        step(m)
    }
}
