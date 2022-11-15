package modify

import (
    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/nfa"
    "mipt_formal/auto/internal/tools"
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

func findUnusedDFS(state common.State, visited []bool, m *nfa.Machine) {
    if visited[state] {
        return
    }
    visited[state] = true

    for _, children := range m.Delta[state] {
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
    newDelta := make([]common.Transition, 0)
    for from := range m.Delta {
        newFrom, found := mask[common.State(from)]
        if !found {
            continue
        }
        for by := range m.Delta[from] {
            for to := range m.Delta[from][by] {
                newTo, ok := mask[to]
                if !ok {
                    continue
                }
                newDelta = append(newDelta, common.Transition{From: newFrom, To: newTo, By: by})
            }
        }
    }

    newStart := make([]common.State, 0)
    for s := range m.Start {
        alias, found := mask[s]
        if !found {
            continue
        }
        newStart = append(newStart, alias)
    }

    newFinal := make([]common.State, 0)
    for f := range m.Final {
        alias, found := mask[f]
        if !found {
            continue
        }
        newFinal = append(newFinal, alias)
    }

    *m = *nfa.NewMachine(newStart, newFinal, newDelta)
}
