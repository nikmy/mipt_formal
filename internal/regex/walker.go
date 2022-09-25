package regex

import (
    "mipt_formal/internal/fsm"
)

func RunDFSWalker(init *IntrusiveState, final *IntrusiveState) ([]fsm.State, []fsm.State, []fsm.Transition) {
    transitions := make([]fsm.Transition, 0)

    walker := &Walker{
        visited: map[*IntrusiveState]bool{
            final: true,
        },
        mapping: map[*IntrusiveState]fsm.State{
            init:  0,
            final: 1,
        },
        current: init,
        lastSID: 1,
    }

    walker.Walk(&transitions)

    return []fsm.State{0}, []fsm.State{1}, transitions
}

type Walker struct {
    visited map[*IntrusiveState]bool
    mapping map[*IntrusiveState]fsm.State
    current *IntrusiveState
    lastSID fsm.State
}

func (w *Walker) Walk(t *[]fsm.Transition) {
    if len(w.current.next) == 0 {
        return
    }

    w.visited[w.current] = true

    cur := w.current

    for _, child := range cur.next {
        if w.visited[child] {
            continue
        }

        id, used := w.mapping[child]
        if !used {
            id = w.lastSID
            w.lastSID++
        }

        *t = append(*t, fsm.Transition{
            From: w.mapping[cur],
            To:   id,
            By:   cur.label,
        })

        w.current = child
        w.Walk(t)
    }

    w.current = cur
}
