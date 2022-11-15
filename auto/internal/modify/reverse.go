package modify

import (
    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/nfa"
)

func Reverse(m *nfa.Machine) {
    invDelta := make([]common.Transition, 0, len(m.Delta))
    for to := range m.Delta {
        for by := range m.Delta[to] {
            for from := range m.Delta[to][by] {
                invDelta = append(invDelta, common.Transition{
                    From: from,
                    To:   common.State(to),
                    By:   by,
                })
            }
        }
    }

    *m = *nfa.NewMachine(m.Final.AsSlice(), m.Start.AsSlice(), invDelta)
}
