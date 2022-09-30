package regex

import (
    "mipt_formal/internal/common"
)

func RunDFSWalker(init *intrusiveState, final *intrusiveState) ([]common.State, []common.State, []common.Transition) {
    transitions := make([]common.Transition, 0)

    walker := &Walker{
        visited: map[*intrusiveState]bool{},
        mapping: map[*intrusiveState]common.State{
            init:  0,
            final: 1,
        },
        current: init,
        lastSID: 1,
    }

    walker.Walk(&transitions)

    return []common.State{0}, []common.State{1}, transitions
}

type Walker struct {
    visited map[*intrusiveState]bool
    mapping map[*intrusiveState]common.State
    current *intrusiveState
    lastSID common.State
}

func (w *Walker) Walk(t *[]common.Transition) {
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

        *t = append(*t, common.Transition{
            From: w.mapping[cur],
            To:   id,
            By:   cur.label,
        })

        w.current = child
        w.Walk(t)
    }

    w.current = cur
}
