package modify

import (
    "testing"

    "github.com/stretchr/testify/assert"

    . "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
)

func TestReverse(t *testing.T) {
    type testcase struct {
        name string
        got  *nfa.Machine
        want *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "just works",
            got: nfa.New([]State{0}, []State{2}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "b"},
            }),
            want: nfa.New([]State{2}, []State{0}, []Transition{
                {From: 2, To: 1, By: "b"},
                {From: 1, To: 0, By: "a"},
            }),
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            Reverse(c.got)
            assert.True(t, c.got.Equal(c.want), c.got.DOA())
        })
    }
}
