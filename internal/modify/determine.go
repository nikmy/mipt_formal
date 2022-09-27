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
           next := make(SetState)
           accept := false
           for from in q:
               for to in from.Outgoing:
                   if from.Accept:
                       next.Accept = true
                   next.push(to)
           if next not in used:
               queue.push(next)
               used.push(next)
               m.NewState(next)

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
    queue := tools.NewQueue[internalState]()
    for _, s := range m.Start {
        accept := false
        for _, f := range m.Final {
            if s == f {
                accept = true
                break
            }
        }
        queue.Push(makeInternal(s, accept, m))
    }

    used := newStateSet(m.NStates())
    accepts := tools.NewSet[common.State](m.Final...)

    for !queue.Empty() {
        q := queue.Pop()
        next := newInternalState(m)
        for from := range q.Mask.Iterate() {
            for _, t := range m.Delta[from] {
                for to := range t {
                    next.Mask.Fix(int(to))
                    if accepts.Has(to) {
                        next.Accept = true
                    }
                }
            }
        }
        if used.TryInsert(next) {
            queue.Push(next)
            addState(m, &next)
        }
    }
}

func addState(m *nfa.Machine, state *internalState) {
    newState := common.State(m.NStates())
    m.Delta = append(m.Delta, make(nfa.Transitions))
    for i := range state.Mask.Iterate() {
        for w, dst := range m.Delta[i] {
            for s := range dst {
                m.AddTransition(newState, s, w)
            }
        }
    }

    if state.Accept {
        m.MarkAsFinal(newState)
    }
}

func newStateSet(initCap int) *stateSet {
    return &stateSet{
        data: make([]internalState, 0, initCap),
    }
}

type stateSet struct {
    data []internalState
}

func (s *stateSet) TryInsert(state internalState) bool {
    for _, m := range s.data {
        if !m.Mask.Xor(state.Mask) {
            return false
        }
    }
    s.data = append(s.data, state)
    return true
}

func makeInternal(state common.State, accept bool, m *nfa.Machine) internalState {
    s := internalState{
        Mask:   tools.NewBitset(m.NStates()),
        Accept: accept,
    }
    s.Mask.Fix(int(state))
    return s
}

func newInternalState(m *nfa.Machine) internalState {
    return internalState{
        Mask:   tools.NewBitset(m.NStates()),
        Accept: false,
    }
}

type internalState struct {
    Mask   *tools.Bitset
    Accept bool
}
