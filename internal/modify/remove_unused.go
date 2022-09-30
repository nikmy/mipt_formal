package modify

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
    "mipt_formal/internal/tools"
)

var _ Modifier = RemoveUnusedStates

func RemoveUnusedStates(stateMachine *nfa.Machine) {
    visited := make([]bool, len(stateMachine.Delta))
    for s := range stateMachine.Start {
        findUnusedDFS(s, visited, stateMachine)
    }

    unused := tools.NewSet[common.State]()
    for i, used := range visited {
        if !used {
            unused.Insert(common.State(i))
        }
    }
    removeUnusedStates(stateMachine, unused)
}

func findUnusedDFS(s common.State, visited []bool, m *nfa.Machine) {
    if visited[s] {
        return
    }
    visited[s] = true

    for _, children := range m.Delta[s] {
        for child := range children {
            findUnusedDFS(child, visited, m)
        }
    }
}

func removeUnusedStates(m *nfa.Machine, toRemove tools.Set[common.State]) {
    mask := calculateRemovingMask(m, toRemove)
    applyMask(m, mask)
}

func calculateRemovingMask(m *nfa.Machine, toRemove tools.Set[common.State]) map[common.State]common.State {
    readPos, writePos := 0, 0
    mask := make(map[common.State]common.State, toRemove.Size()) // old --> new

    // hack for break external cycle
    func() {
        for readPos < len(m.Delta) {
            for toRemove.Has(common.State(readPos)) {
                if readPos == len(m.Delta)-1 {
                    return
                }
                readPos++
            }
            mask[common.State(readPos)] = common.State(writePos)
            writePos++
            readPos++
        }
    }()

    return mask
}

func applyMask(m *nfa.Machine, mask map[common.State]common.State) {
    for from := range m.Delta {
        for _, to := range m.Delta[from] {
            for x := range to {
                if alias, found := mask[x]; found && alias != x {
                    to.Delete(x)
                    to.Insert(mask[x])
                }
            }
        }
    }

    newDelta := make([]nfa.Transitions, len(mask))
    for i := range newDelta {
        newDelta[i] = make(nfa.Transitions, len(m.Delta[i]))
    }
    for i := range m.Delta {
        alias, found := mask[common.State(i)]
        if !found {
            continue
        }
        newDelta[alias] = m.Delta[i]
    }
    m.Delta = newDelta

    newStart := tools.NewSet[common.State]()
    for s := range m.Start {
        alias, found := mask[s]
        if !found {
            continue
        }
        newStart.Insert(alias)
    }
    m.Start = newStart

    newFinal := tools.NewSet[common.State]()
    for f := range m.Final {
        alias, found := mask[f]
        if !found {
            continue
        }
        newFinal.Insert(alias)
    }
    m.Final = newFinal
}
