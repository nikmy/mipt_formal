package modify

import (
    "testing"

    "github.com/stretchr/testify/assert"

    . "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/nfa"
)

func TestComplete(t *testing.T) {
    type testcase struct {
        name string
        abc  string
        got  *nfa.Machine
        want *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "fantom state",
            abc:  "a",
            got:  nfa.NewMachine([]State{0}, []State{0}, []Transition{}),
            want: nfa.NewMachine([]State{0}, []State{0}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 1, By: "a"},
            }),
        },
        {
            name: "just works",
            abc:  "ab",
            got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "b"},
                {From: 2, To: 1, By: "a"},
            }),
            want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 3, By: "b"},
                {From: 1, To: 3, By: "a"},
                {From: 1, To: 2, By: "b"},
                {From: 2, To: 1, By: "a"},
                {From: 2, To: 3, By: "b"},
                {From: 3, To: 3, By: "a"},
                {From: 3, To: 3, By: "b"},
            }),
        },
        {
            name: "already complete",
            abc:  "ab",
            got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 3, By: "b"},
                {From: 1, To: 3, By: "a"},
                {From: 1, To: 2, By: "b"},
                {From: 2, To: 1, By: "a"},
                {From: 2, To: 1, By: "b"},
                {From: 3, To: 3, By: "a"},
                {From: 3, To: 3, By: "b"},
            }),
            want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 3, By: "b"},
                {From: 1, To: 3, By: "a"},
                {From: 1, To: 2, By: "b"},
                {From: 2, To: 1, By: "a"},
                {From: 2, To: 1, By: "b"},
                {From: 3, To: 3, By: "a"},
                {From: 3, To: 3, By: "b"},
            }),
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            Complete(c.abc)(c.got)
            assert.True(t, c.got.Equal(c.want))
        })
    }
}
