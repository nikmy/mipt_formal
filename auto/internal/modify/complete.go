package modify

import (
	"mipt_formal/internal/common"
	"mipt_formal/internal/nfa"
)

func Complete(alphabet string) Modifier {
	return func(m *nfa.Machine) {
		fantomState := common.State(m.NStates())
		useFantom := false
		for from := range m.Delta {
			for _, symbol := range alphabet {
				by := common.Word(symbol)
				if _, has := m.Delta[from][by]; !has {
					m.AddTransition(common.State(from), fantomState, by)
					useFantom = true
				}
			}
		}
		if useFantom {
			m.Delta = append(m.Delta, make(nfa.Transitions, len(alphabet)))
			for _, symbol := range alphabet {
				m.AddTransition(fantomState, fantomState, common.Word(symbol))
			}
		}
	}
}
