package modify

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
)

func Complete(alphabet string) Modifier {
    return func(m *nfa.Machine) {
        fantomState := common.State(m.NStates())
        useFantom := false
        for s, t := range m.Delta {
            for _, c := range alphabet {
                w := common.Word(c)
                if _, has := t[w]; !has {
                    m.AddTransition(common.State(s), fantomState, w)
                    useFantom = true
                }
            }
        }
        if useFantom {
            m.Delta = append(m.Delta, make(nfa.Transitions, len(alphabet)))
            for _, c := range alphabet {
                m.AddTransition(fantomState, fantomState, common.Word(c))
            }
        }
    }
}
