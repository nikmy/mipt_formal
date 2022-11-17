package modify

import (
    "github.com/stretchr/testify/assert"
    "mipt_formal/gram/internal/common"
    "testing"

    "mipt_formal/gram/internal/cf"
)

func TestChomskyNormalizer_Run(t *testing.T) {
    type testcase struct {
        name  string
        rules []string
    }

    cases := [...]testcase{
        {
            name: "already in CNF",
            rules: []string{
                "S -> _",
                "S -> AB",
                "A -> AA",
                "A -> a",
                "B -> b",
            },
        },
        {
            name: "long rule",
            rules: []string{
                "S -> ABC",
                "A -> a",
                "B -> b",
                "C -> c",
            },
        },
        {
            name: "null-productive",
            rules: []string{
                "S -> A",
                "A -> _",
            },
        },
        {
            name: "non-generative",
            rules: []string{
                "S -> aA",
            },
        },
        {
            name: "unit-rules chain",
            rules: []string{
                "S -> A", "A -> B", "B -> C", "C -> c",
            },
        },
        {
            name: "mixed rule",
            rules: []string{
                "S -> aSbS",
                "S -> _",
            },
        },
        {
            name: "unreachable non-terminals",
            rules: []string{
                "S -> A",
                "A -> a",
                "B -> C",
                "C -> c",
            },
        },
        {
            name: "null rule",
            rules: []string{
                "S -> _",
            },
        },
        {
            name: "just works",
            rules: []string{
                "S -> ABcS", "S -> _",
                "A -> AbbAB", "A -> SA", "A -> a",
                "B -> b", "B -> _",
            },
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            g, _ := cf.ParseGrammar(common.NewStringsReader(tc.rules))
            normalizer := NewChomskyNormalizer(g)
            normalizer.Run()
            assert.True(t, normalizer.checkNF())
        })
    }
}

func TestChomskyNormalizer_checkNF(t *testing.T) {
    type testcase struct {
        name  string
        rules []string
        want  bool // false by default
    }

    cases := [...]testcase{
        {
            name: "already in CNF",
            rules: []string{
                "S -> _",
                "S -> AB",
                "A -> AA",
                "A -> a",
                "B -> b",
            },
            want: true,
        },
        {
            name: "long rule",
            rules: []string{
                "S -> ABC",
                "A -> a",
                "B -> b",
                "C -> c",
            },
        },
        {
            name: "null-productive",
            rules: []string{
                "S -> A",
                "A -> _",
            },
        },
        {
            name: "non-generative",
            rules: []string{
                "S -> aA",
            },
        },
        {
            name: "long unit-rules chain",
            rules: []string{
                "S -> A", "A -> B", "B -> C", "C -> D", "D -> E", "E -> F", "F -> G", "G -> H",
                "H -> I", "I -> J", "J -> K", "K -> L", "L -> M", "M -> N", "N -> O", "O -> P",
                "P -> Q", "Q -> R", "R -> T", "T -> U", "U -> V", "V -> W", "W -> X", "X -> Y",
                "Y -> Z", "Z -> a",
            },
        },
        {
            name: "mixed rule",
            rules: []string{
                "S -> aSbS",
                "S -> _",
            },
        },
        {
            name: "unreachable non-terminals",
            rules: []string{
                "S -> A",
                "A -> a",
                "B -> C",
                "C -> c",
            },
        },
        {
            name: "null rule",
            rules: []string{
                "S -> _",
            },
            want: true,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            g, _ := cf.ParseGrammar(common.NewStringsReader(tc.rules))
            assert.Equal(t, tc.want, NewChomskyNormalizer(g).checkNF())
        })
    }
}
