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
    clearFinals(m, toRemove, mask)
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
            if readPos > writePos {
                mask[common.State(readPos)] = common.State(writePos)
                m.Delta[writePos] = m.Delta[readPos]
            }
            writePos++
            readPos++
        }
    }()

    m.Delta = m.Delta[:writePos]
    return mask
}

func clearFinals(m *nfa.Machine, toRemove tools.Set[common.State], mask map[common.State]common.State) {
    newFinals := tools.NewSet[common.State]()
    for f := range m.Final {
        if toRemove.Has(f) {
            continue
        }
        if alias, found := mask[f]; found {
            newFinals.Insert(alias)
        } else {
            newFinals.Insert(f)
        }
    }
    m.Final = newFinals
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
}
