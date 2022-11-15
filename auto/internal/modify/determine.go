package modify

import (
    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/nfa"
    "mipt_formal/auto/internal/tools"
)

// Determine build DFA, equal to the given NFA.
// !!! You need to remove unused states after it explicitly !!!
/*
   Algorithm:
       m       : given NFA
       q_start : NFA start states

       queue := [q_start]
       used := Set[SetState]
       while queue not empty:
           q := queue.pop()
           next := make(map[Word]SetState)
           for from in q:
               for by, to in from.Outgoing:
                   if from.Accept:
                       next[by].Accept = true
                   next.push(to)
           for state in next.Values():
               if state not in used:
                   queue.push(state)
                   used.push(state)
                   m.NewState(state)
        rename_states()

   Heuristics:
       SetState := BitSet<m.NStates>
       used.Has := AllOf({ u ^ x for u in used})

   Complexity:
       n := |Q|
       k := |A| - constant
       T(used.Has) = O(n / 64)
       T(for from in q) = O(n * k)
       T(all) := O(2^n * (n * k + n * 2^n / 64)) = O(n * 2^n), const = |A| * 2^(-n) + 1 / 64
*/
func Determine(m *nfa.Machine) {
    queue := tools.NewQueue[*internalState]()

    used := newStateSet(m.NStates())
    aliases := make(map[*internalState]common.State)

    nStates := m.NStates()

    start := newInternalState(nStates)
    for state := range m.Start {
        if m.Final.Has(state) {
            start.Accept = true
        }
        start.Mask.Fix(int(state))
    }
    if start.Mask.Count() > 1 {
        newStart := addState(m, start)
        aliases[start] = newStart
        m.Start = tools.NewSet[common.State](newStart)
    }

    used.TryInsert(start)
    queue.Push(start)

    for !queue.Empty() {
        state := queue.Pop()
        group := make(map[common.Word]*internalState, 0)
        for from := range state.Mask.Iterate() {
            for by := range m.Delta[from] {
                if _, found := group[by]; !found {
                    group[by] = newInternalState(nStates)
                }
                for to := range m.Delta[from][by] {
                    group[by].Mask.Fix(int(to))
                    if m.Final.Has(to) {
                        group[by].Accept = true
                    }
                }
            }
        }
        for _, next := range group {
            if used.TryInsert(next) {
                queue.Push(next)
                if next.Mask.Count() > 1 {
                    aliases[next] = addState(m, next)
                }
            }
        }
    }

    for i := range m.Delta { // update ingoing edges to new states
        bySym := make(map[common.Word]*internalState, 0)
        for word := range m.Delta[i] {
            if _, found := bySym[word]; !found {
                bySym[word] = newInternalState(nStates)
            }
            for state := range m.Delta[i][word] {
                bySym[word].Mask.Fix(int(state))
            }
        }
        for word, state := range bySym {
            if state.Mask.Count() <= 1 { // skip old states
                continue
            }
            p := used.Find(state)
            if p == nil { // skip unused states
                continue
            }
            newTo := aliases[p]
            m.Delta[i][word] = tools.NewSet[common.State](newTo)
        }
    }
}

func addState(m *nfa.Machine, state *internalState) common.State {
    newState := m.AddState()
    for from := range state.Mask.Iterate() { // add outgoing edges from new state
        for word := range m.Delta[from] {
            for to := range m.Delta[from][word] {
                m.AddTransition(newState, to, word)
            }
        }
    }

    if state.Accept {
        m.MarkAsFinal(newState)
    }

    return newState
}

func newStateSet(initCap int) *stateSet {
    return &stateSet{
        data: make([]*internalState, 0, initCap),
    }
}

type stateSet struct {
    data []*internalState
}

func (s *stateSet) TryInsert(state *internalState) bool {
    for _, m := range s.data {
        if !m.Mask.Xor(state.Mask) {
            return false
        }
    }
    s.data = append(s.data, state)
    return true
}

func (s *stateSet) Find(state *internalState) *internalState {
    for _, m := range s.data {
        if !m.Mask.Xor(state.Mask) {
            return m
        }
    }
    return nil
}

func newInternalState(nStates int) *internalState {
    return &internalState{
        Mask:   tools.NewBitset(nStates),
        Accept: false,
    }
}

type internalState struct {
    Mask   *tools.Bitset
    Accept bool
}
