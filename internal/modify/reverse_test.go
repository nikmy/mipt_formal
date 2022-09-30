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

func Test_singleStart(t *testing.T) {
    type testcase struct {
        name string
        got  *nfa.Machine
        want *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "already single start",
            got: nfa.New([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: Epsilon},
            }),
            want: nfa.New([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: Epsilon},
            }),
        },
        {
            name: "just works",
            got: nfa.New([]State{0, 1, 2}, []State{3}, []Transition{
                {From: 0, To: 3, By: Epsilon},
                {From: 1, To: 3, By: Epsilon},
                {From: 2, To: 3, By: Epsilon},
            }),
            want: nfa.New([]State{1}, []State{0}, []Transition{
                {From: 1, To: 0, By: Epsilon},
            }),
        },
        {
            name: "correct acceptances",
            got: nfa.New([]State{0, 1, 2}, []State{2, 3}, []Transition{
                {From: 0, To: 3, By: Epsilon},
                {From: 1, To: 3, By: Epsilon},
                {From: 2, To: 3, By: Epsilon},
            }),
            want: nfa.New([]State{1}, []State{0, 1}, []Transition{
                {From: 1, To: 0, By: Epsilon},
            }),
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            singleStart(c.got)
            assert.True(t, c.got.Equal(c.want))
        })
    }
}

func Test_singleAcceptance(t *testing.T) {
    type testcase struct {
        name string
        got  *nfa.Machine
        want *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "already single acceptance",
            got: nfa.New([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: Epsilon},
            }),
            want: nfa.New([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: Epsilon},
            }),
        },
        {
            name: "just works",
            got: nfa.New([]State{0}, []State{1, 2, 3}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 2, By: "b"},
                {From: 0, To: 3, By: "c"},
            }),
            want: nfa.New([]State{0}, []State{1}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 1, By: "b"},
                {From: 0, To: 1, By: "c"},
            }),
        },
        {
            name: "correct start",
            got: nfa.New([]State{0}, []State{0, 1, 2, 3}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 0, To: 2, By: "b"},
                {From: 0, To: 3, By: "c"},
            }),
            want: nfa.New([]State{0}, []State{0}, []Transition{
                {From: 0, To: 0, By: "a"},
                {From: 0, To: 0, By: "b"},
                {From: 0, To: 0, By: "c"},
            }),
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            singleAcceptance(c.got)
            assert.True(t, c.got.Equal(c.want))
        })
    }
}
