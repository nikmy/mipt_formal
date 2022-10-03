package modify

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
    "mipt_formal/internal/tools"
)

// Determine build DFA, equal to the given NFA.
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

    start := newInternalState(m)
    for s := range m.Start {
        if m.Final.Has(s) {
            start.Accept = true
        }
        start.Mask.Fix(int(s))
    }
    if start.Mask.Count() > 1 {
        newStart := addState(m, start)
        aliases[start] = newStart
        m.Start = tools.NewSet[common.State](newStart)
    }

    used.TryInsert(start)
    queue.Push(start)

    for !queue.Empty() {
        q := queue.Pop()
        group := make(map[common.Word]*internalState, 0)
        for from := range q.Mask.Iterate() {
            for by, t := range m.Delta[from] {
                if _, found := group[by]; !found {
                    group[by] = newInternalState(m)
                }
                for to := range t {
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
                    if next.Accept {
                        m.Final.Insert(aliases[next])
                    }
                }
            }
        }
    }

    for i, transitions := range m.Delta { // update ingoing edges to new states
        bySym := make(map[common.Word]*internalState, 0)
        for w, dst := range transitions {
            if _, found := bySym[w]; !found {
                bySym[w] = newInternalState(m)
            }
            for s := range dst {
                bySym[w].Mask.Fix(int(s))
            }
        }
        for w, dst := range transitions {
            if dst.Size() > 1 {
                delete(m.Delta[i], w)
            }
        }
        for w, s := range bySym {
            if s.Mask.Count() <= 1 { // skip old states
                continue
            }
            p := used.Find(s)
            if p == nil { // skip unused states
                continue
            }
            newTo := aliases[p]
            m.Delta[i][w] = tools.NewSet[common.State](newTo)
        }
    }
}

func addState(m *nfa.Machine, state *internalState) common.State {
    newState := m.AddState()
    for i := range state.Mask.Iterate() { // add outgoing edges from new state
        for w, dst := range m.Delta[i] {
            for s := range dst {
                m.AddTransition(newState, s, w)
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

func newInternalState(m *nfa.Machine) *internalState {
    return &internalState{
        Mask:   tools.NewBitset(m.NStates()),
        Accept: false,
    }
}

type internalState struct {
    Mask   *tools.Bitset
    Accept bool
}
