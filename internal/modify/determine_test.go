package modify

import (
    "testing"

    "github.com/stretchr/testify/assert"

    . "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
)

func TestDetermine(t *testing.T) {
    type testcase struct {
        name string
        got  *nfa.Machine
        want *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "just works 1",
            got: nfa.New([]State{0}, []State{1, 3}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 2, By: "a"},
                {From: 2, To: 3, By: "a"},
            }),
            want: nfa.New([]State{0}, []State{1, 3, 4}, []Transition{
                {From: 0, To: 4, By: "a"},
                {From: 2, To: 3, By: "a"},
                {From: 4, To: 3, By: "a"},
            }),
        },
        {
            name: "just works 2",
            got: nfa.New([]State{0}, []State{4}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "a"},
                {From: 1, To: 3, By: "a"},
                {From: 2, To: 4, By: "a"},
                {From: 3, To: 4, By: "a"},
            }),
            want: nfa.New([]State{0}, []State{4}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 5, By: "a"},
                {From: 5, To: 4, By: "a"},
                {From: 2, To: 4, By: "a"},
                {From: 3, To: 4, By: "a"},
            }),
        },
        {
            name: "just works 3",
            got: nfa.New([]State{0}, []State{0, 1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 0, By: "a"},
            }),
            want: nfa.New([]State{0}, []State{0, 1, 2}, []Transition{
                {From: 0, To: 2, By: "a"},
                {From: 2, To: 2, By: "a"},
            }),
        },
        {
            name: "two starts",
            got: nfa.New([]State{0, 1}, []State{2}, []Transition{
                {From: 0, To: 2, By: "a"},
                {From: 1, To: 2, By: "a"},
            }),
            want: nfa.New([]State{3}, []State{2}, []Transition{
                {From: 0, To: 2, By: "a"},
                {From: 1, To: 2, By: "a"},
                {From: 3, To: 2, By: "a"},
            }),
        },
        {
            name: "already determined",
            got: nfa.New([]State{0}, []State{3}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "a"},
                {From: 2, To: 3, By: "a"},
            }),
            want: nfa.New([]State{0}, []State{3}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 2, By: "a"},
                {From: 2, To: 3, By: "a"},
            }),
        },
        {
            name: "ignore unused states & edges",
            got: nfa.New([]State{0}, []State{4}, []Transition{
                {From: 0, To: 4, By: "b"},
                {From: 1, To: 2, By: "a"},
                {From: 1, To: 3, By: "a"},
            }),
            want: nfa.New([]State{0}, []State{4}, []Transition{
                {From: 0, To: 4, By: "b"},
                {From: 1, To: 2, By: "a"},
                {From: 1, To: 3, By: "a"},
            }),
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            Determine(c.got)
            assert.True(t, c.got.Equal(c.want), c.got.DOA())
        })
    }
}
