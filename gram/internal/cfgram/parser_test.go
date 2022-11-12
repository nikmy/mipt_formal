package cfgram

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "io"
    "testing"
)

func newMockReader(s string) *mockReader {
    lines := bytes.Split([]byte(s), []byte("\n"))
    return &mockReader{
        lines: lines,
        curr:  0,
    }
}

type mockReader struct {
    lines [][]byte
    curr  int
}

func (r *mockReader) ReadLine() ([]byte, error) {
    if r.curr == len(r.lines) {
        return nil, io.EOF
    }
    next := r.lines[r.curr]
    r.curr++
    return next, nil
}

func TestParseGrammar(t *testing.T) {
    type testcase struct {
        name  string
        arg   string
        want  *Grammar
        throw bool
    }

    cases := [...]testcase{
        {
            name: "no start symbol",
            arg:  "A -> aA\n",
            want: nil,
        },
        {
            name: "not context-free",
            arg: "S -> _\n" +
                    "aA -> b",
            want: nil,
        },
        {
            name: "not correct format",
            arg:  "S = _\n",
            want: nil,
        },
        {
            name: "just works",
            arg: "S -> aS\n" +
                    "S -> _",
            want: &Grammar{
                Rules: []Rule{
                    {
                        Left:  NonTerminal('S'),
                        Right: []Symbol{Symbol('a'), Symbol('S')},
                    },
                    {
                        Left:  NonTerminal('S'),
                        Right: []Symbol{Epsilon},
                    },
                }},
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            r := newMockReader(tc.arg)
            got, _ := ParseGrammar(r)
            assert.Equal(t, tc.want, got)
        })
    }
}
