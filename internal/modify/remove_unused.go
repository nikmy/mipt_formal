package modify

import (
    "mipt_formal/internal/nfa"
    "mipt_formal/internal/tools"
)

var _ Modifier = RemoveUnusedStates

func RemoveUnusedStates(m *nfa.Machine) {
    visited := make([]bool, len(m.Delta))
    findUnusedDFS(m.Start[0], visited, m)

    unused := make([]nfa.State, 0)
    for i, used := range visited {
        if !used {
            unused = append(unused, nfa.State(i))
        }
    }
    removeStates(m, unused)
}

func findUnusedDFS(s nfa.State, visited []bool, m *nfa.Machine) {
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

func removeStates(m *nfa.Machine, states []nfa.State) {
    toRemove := tools.NewSet[nfa.State](states...)

    mask := calculateRemovingMask(m, toRemove)

    removeOutgoingEdges(m, toRemove)

    applyMask(m, mask)

    clearFinals(m, toRemove, mask)
}

func removeOutgoingEdges(m *nfa.Machine, states tools.Set[nfa.State]) {
    for from := range m.Delta {
        if states.Has(nfa.State(from)) {
            continue
        }
        for w, to := range m.Delta[from] {
            for x := range to {
                if states.Has(x) {
                    m.Delta[from][w].Delete(x)
                }
            }
        }
    }
}

func clearFinals(m *nfa.Machine, toRemove tools.Set[nfa.State], mask map[nfa.State]nfa.State) {
    readPos, writePos := 0, 0
    for readPos < len(m.Final) {
        if toRemove.Has(m.Final[readPos]) {
            if readPos == len(m.Final)-1 {
                break
            }
            readPos++
        }
        m.Final[writePos] = mask[m.Final[readPos]]
        writePos++
        readPos++
    }
    m.Final = m.Final[:writePos]
}

func calculateRemovingMask(m *nfa.Machine, toRemove tools.Set[nfa.State]) map[nfa.State]nfa.State {
    readPos, writePos := 0, 0
    mask := make(map[nfa.State]nfa.State, toRemove.Size()) // old --> new

    // hack for break external cycle
    func() {
        for readPos < len(m.Delta) {
            for toRemove.Has(nfa.State(readPos)) {
                if readPos == len(m.Delta)-1 {
                    return
                }
                readPos++
            }
            if readPos > writePos {
                mask[nfa.State(readPos)] = nfa.State(writePos)
                m.Delta[writePos] = m.Delta[readPos]
            }
            writePos++
            readPos++
        }
    }()

    m.Delta = m.Delta[:writePos]
    return mask
}

func applyMask(m *nfa.Machine, mask map[nfa.State]nfa.State) {
    for from := range m.Delta {
        for _, to := range m.Delta[from] {
            for x := range to {
                if mask[x] != x {
                    to.Delete(x)
                    to.Insert(mask[x])
                }
            }
        }
    }
}
