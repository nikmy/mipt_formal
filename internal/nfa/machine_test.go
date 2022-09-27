package nfa

import (
    "fmt"
    "mipt_formal/internal/common"
    "mipt_formal/internal/doa"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestFiniteStateMachine_DOA(t *testing.T) {

    // many states cannot be tested, because unordered map is used

    type testcase struct {
        name     string
        machine  Machine
        expected string
    }

    cases := [...]testcase{
        {
            name: "the only state",
            machine: Machine{
                Delta: []transitions{0: nil},
                Start: []common.State{0},
                Final: []common.State{0},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 0) +
                    fmt.Sprintf(doa.AcceptanceFormat, 0) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 0) + doa.End,
        },
        {
            name: "one edge",
            machine: Machine{
                Delta: []transitions{
                    0: map[common.Word][]common.State{"a": {0}},
                },
                Start: []common.State{0},
                Final: []common.State{0},
            },
            expected: doa.Version + fmt.Sprintf(doa.StartFormat, 0) +
                    fmt.Sprintf(doa.AcceptanceFormat, 0) + doa.Begin +
                    fmt.Sprintf(doa.StateFormat, 0) +
                    fmt.Sprintf(doa.EdgeFormat, "a", 0) + doa.End,
        },
        {
            name: "epsilon transition",
            machine: Machine{
                Delta: []transitions{
                    0: map[common.Word][]common.State{"": {0}},
                },
                Start: []common.State{0},
                Final: []common.State{0},
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
