package regex

import (
    "mipt_formal/internal/nfa"
)

func RunDFSWalker(init *IntrusiveState, final *IntrusiveState) ([]nfa.State, []nfa.State, []nfa.Transition) {
    transitions := make([]nfa.Transition, 0)

    walker := &Walker{
        visited: map[*IntrusiveState]bool{},
        mapping: map[*IntrusiveState]nfa.State{
            init:  0,
            final: 1,
        },
        current: init,
        lastSID: 1,
    }

    walker.Walk(&transitions)

    return []nfa.State{0}, []nfa.State{1}, transitions
}

type Walker struct {
    visited map[*IntrusiveState]bool
    mapping map[*IntrusiveState]nfa.State
    current *IntrusiveState
    lastSID nfa.State
}

func (w *Walker) Walk(t *[]nfa.Transition) {
    if len(w.current.next) == 0 {
        return
    }

    if w.visited[w.current] {
        return
    }

    w.visited[w.current] = true

    cur := w.current

    for _, child := range cur.next {
        id, used := w.mapping[child]
        if !used {
            w.lastSID++
            id = w.lastSID
            w.mapping[child] = id
        }

        *t = append(*t, nfa.Transition{
            From: w.mapping[cur],
            To:   id,
            By:   cur.label,
        })

        w.current = child
        w.Walk(t)
    }

    w.current = cur
}
