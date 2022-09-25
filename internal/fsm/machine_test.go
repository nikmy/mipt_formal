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
        machine  NFA
        expected string
    }

    cases := [...]testcase{
        {
            name: "the only state",
            machine: NFA{
                delta: []transitions{0: nil},
                start: []State{0},
                final: []State{0},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 0) +
                    fmt.Sprintf(doa.AcceptanceFormat, 0) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 0) + doa.End,
        },
        {
            name: "one edge",
            machine: NFA{
                delta: []transitions{
                    0: map[Word][]State{"a": {0}},
                },
                start: []State{0},
                final: []State{0},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 0) +
                    fmt.Sprintf(doa.AcceptanceFormat, 0) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 0) +
                    fmt.Sprintf(doa.EdgeFormat, "a", 0) + doa.End,
        },
        {
            name: "epsilon transition",
            machine: NFA{
                delta: []transitions{
                    0: map[Word][]State{"": {0}},
                },
                start: []State{0},
                final: []State{0},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 0) +
                    fmt.Sprintf(doa.AcceptanceFormat, 0) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 0) +
                    fmt.Sprintf(doa.EdgeFormat, doa.Epsilon, 0) + doa.End,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            assert.Equal(t, c.expected, c.machine.DOA())
        })
    }
}
