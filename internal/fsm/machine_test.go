package fsm

import (
    "fmt"
    "mipt_formal/internal/doa"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestFiniteStateMachine_DOA(t *testing.T) {

    // many states cannot be tested, because unordered map is used

    type testcase struct {
        name     string
        machine  finiteStateMachine
        expected string
    }

    cases := [...]testcase{
        {
            name: "the only state",
            machine: finiteStateMachine{
                delta: map[State]transitions{1: nil},
                start: []State{1},
                final: []State{1},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 1) +
                    fmt.Sprintf(doa.AcceptanceFormat, 1) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 1) + doa.End,
        },
        {
            name: "one edge",
            machine: finiteStateMachine{
                delta: map[State]transitions{
                    1: map[Word]State{"a": 1},
                },
                start: []State{1},
                final: []State{1},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 1) +
                    fmt.Sprintf(doa.AcceptanceFormat, 1) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 1) +
                    fmt.Sprintf(doa.EdgeFormat, "a", 1) + doa.End,
        },
        {
            name: "epsilon transition",
            machine: finiteStateMachine{
                delta: map[State]transitions{
                    1: map[Word]State{"": 1},
                },
                start: []State{1},
                final: []State{1},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 1) +
                    fmt.Sprintf(doa.AcceptanceFormat, 1) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 1) +
                    fmt.Sprintf(doa.EdgeFormat, doa.Epsilon, 1) + doa.End,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            assert.Equal(t, c.expected, c.machine.DOA())
        })
    }
}
