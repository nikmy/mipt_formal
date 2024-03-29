package modify

import (
    "testing"

    "github.com/stretchr/testify/assert"

    . "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/nfa"
)

func TestRemoveUnusedStates(t *testing.T) {
    type testcase struct {
        name string
        got  *nfa.Machine
        want *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "just works",
            got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 2, To: 1, By: "b"},
                {From: 3, To: 1, By: "c"},
            }),
            want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
            }),
        },
        {
            name: "remove extra acceptances",
            got:  nfa.NewMachine([]State{0}, []State{0, 1}, []Transition{}),
            want: nfa.NewMachine([]State{0}, []State{0}, []Transition{}),
        },
        {
            name: "loop",
            got: nfa.NewMachine([]State{0}, []State{1, 3}, []Transition{
                {From: 0, To: 2, By: "a"},
                {From: 2, To: 0, By: "b"},
                {From: 0, To: 1, By: "c"},
                {From: 3, To: 1, By: "a"},
            }),
            want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 2, By: "a"},
                {From: 2, To: 0, By: "b"},
                {From: 0, To: 1, By: "c"},
            }),
        },
        {
            name: "no unused",
            got: nfa.NewMachine([]State{0}, []State{2}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "b"},
            }),
            want: nfa.NewMachine([]State{0}, []State{2}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "b"},
            }),
        },
        {
            name: "isolated",
            got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 2, To: 3, By: "b"},
                {From: 3, To: 2, By: "b"},
            }),
            want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
            }),
        },
        {
            name: "edges to unused",
            got: nfa.NewMachine([]State{0}, []State{2, 6, 3, 7}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 6, By: "b"},
                {From: 6, To: 3, By: "c"},
            }),
            want: nfa.NewMachine([]State{0}, []State{3, 2}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 3, By: "b"},
                {From: 3, To: 2, By: "c"},
            }),
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            RemoveUnusedStates(c.got)
            assert.True(t, c.got.Equal(c.want))
        })
    }
}
