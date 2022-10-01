package modify

import (
    "testing"

    "github.com/stretchr/testify/assert"

    . "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
)

type mockLogger struct{}

func (e mockLogger) Fatal(err error)                  {}
func (e mockLogger) Info(msg string)                  {}
func (e mockLogger) Infof(format string, args ...any) {}

func TestSequence(t *testing.T) {
    type testcase struct {
        name  string
        steps []Step
        init  *nfa.Machine
        want  *nfa.Machine
    }

    cases := [...]testcase{
        {
            name: "just works",
            init: nfa.New([]State{}, []State{}, []Transition{}),
            steps: []Step{
                {
                    Name: "",
                    Func: func(machine *nfa.Machine) {
                        *machine = *nfa.New([]State{0}, []State{0}, []Transition{
                            {From: 0, To: 1, By: "a"},
                            {From: 1, To: 1, By: "a"},
                        })
                    },
                },
            },
            want: nfa.New([]State{0}, []State{0}, []Transition{
                {From: 0, To: 1, By: "a"},
                {From: 1, To: 1, By: "a"},
            }),
        },
        {
            name: "nil machine",
            init: nil,
            steps: []Step{
                {
                    Name: "",
                    Func: func(machine *nfa.Machine) {
                        *machine = *nfa.New([]State{0}, []State{0}, []Transition{
                            {From: 0, To: 1, By: "a"},
                            {From: 1, To: 1, By: "a"},
                        })
                    },
                },
            },
            want: nil,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            Sequence(mockLogger{}, c.steps...)(c.init)
            assert.True(t, c.init.Equal(c.want))
        })
    }
}
