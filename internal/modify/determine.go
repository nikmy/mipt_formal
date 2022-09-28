package modify

import (
    "fmt"
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
           accept := false
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
    start := newInternalState(m)
    for _, s := range m.Start {
        for _, f := range m.Final {
            if s == f {
                start.Accept = true
                break
            }
        }
        start.Mask.Fix(int(s))
    }
    queue.Push(start)

    accepts := tools.NewSet[common.State](m.Final...)

    used := newStateSet(m.NStates())
    aliases := make(map[*internalState]common.State)

    if used.TryInsert(start) && start.Mask.Count() > 1 {
        aliases[start] = addState(m, start)
    }

    for !queue.Empty() {
        q := queue.Pop()
        group := make(map[common.Word]*internalState, 0)
        for from := range q.Mask.Iterate() {
            for w, t := range m.Delta[from] {
                if _, found := group[w]; !found {
                    group[w] = newInternalState(m)
                }
                for to := range t {
                    group[w].Mask.Fix(int(to))
                    if accepts.Has(to) && !group[w].Accept {
                        group[w].Accept = true
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
                        accepts.Insert(aliases[next])
                    }
                }
            }
        }
    }

    for s, a := range aliases {
        fmt.Print("{ ")
        for x := range s.Mask.Iterate() {
            fmt.Printf("%v ", x)
        }
        fmt.Printf("} --> %v\n", a)
    }

    for i, transitions := range m.Delta {
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
            if s.Mask.Count() <= 1 {
                continue
            }
            newTo := aliases[used.Find(s)]
            m.Delta[i][w] = tools.NewSet[common.State](newTo)
        }
    }
}

func addState(m *nfa.Machine, state *internalState) common.State {
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
    panic("wrong use!")
}

func newInternalState(m *nfa.Machine) *internalState {
    return &internalState{
        Mask:   tools.NewBitset(m.NStates()),
        Accept: false,
    }
}

func makeInternal(state common.State, accept bool, m *nfa.Machine) *internalState {
    s := internalState{
        Mask:   tools.NewBitset(m.NStates()),
        Accept: accept,
    }
    s.Mask.Fix(int(state))
    return &s
}

type internalState struct {
    Mask   *tools.Bitset
    Accept bool
}
