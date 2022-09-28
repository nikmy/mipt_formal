package modify

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
    "mipt_formal/internal/tools"
)

var _ Modifier = RemoveUnusedStates

func RemoveUnusedStates(stateMachine *nfa.Machine) {
    visited := make([]bool, len(stateMachine.Delta))
    findUnusedDFS(stateMachine.Start[0], visited, stateMachine)

    unused := make([]common.State, 0)
    for i, used := range visited {
        if !used {
            unused = append(unused, common.State(i))
        }
    }
    removeStates(stateMachine, unused)
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

func removeStates(m *nfa.Machine, states []common.State) {
    toRemove := tools.NewSet[common.State](states...)

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
    readPos, writePos := 0, 0
    for readPos < len(m.Final) {
        if toRemove.Has(m.Final[readPos]) {
            if readPos == len(m.Final)-1 {
                break
            }
            readPos++
        }
        if alias, found := mask[m.Final[readPos]]; found {
            m.Final[writePos] = alias
        } else {
            m.Final[writePos] = m.Final[readPos]
        }
        writePos++
        readPos++
    }
    m.Final = m.Final[:writePos]
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
