package modify

import (
    "github.com/stretchr/testify/assert"
    "mipt_formal/gram/internal/common"
    "mipt_formal/tools"
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
            assert.True(t, normalizer.isNormalForm())
        })
    }
}

func TestChomskyNormalizer_removeNonGenerative(t *testing.T) {
    t.Run("Non-Generative rule", func(t *testing.T) {
        g := cf.NewGrammar([]cf.Rule{
            {
                Left:  cf.Start,
                Right: "aA",
            },
            {
                Left:  cf.Start,
                Right: "",
            },
        })
        NewChomskyNormalizer(g).removeNonGenerative()
        want := []cf.Rule{
            {
                Left:  cf.Start,
                Right: "",
            },
        }
        assert.Equal(t, want, g.Rules)
    })
}

func TestChomskyNormalizer_removeUnreachable(t *testing.T) {
    t.Run("Unreachable non-terminals", func(t *testing.T) {
        g := cf.NewGrammar([]cf.Rule{
            {
                Left:  cf.Start,
                Right: "A",
            },
            {
                Left:  'A',
                Right: "a",
            },
            {
                Left:  'B',
                Right: "C",
            },
            {
                Left:  'C',
                Right: "c",
            },
        })
        NewChomskyNormalizer(g).removeUnreachable()
        want := []cf.Rule{
            {
                Left:  cf.Start,
                Right: "A",
            },
            {
                Left:  'A',
                Right: "a",
            },
        }
        assert.Equal(t, want, g.Rules)
    })
}

func TestChomskyNormalizer_removeMixed(t *testing.T) {
    t.Run("mixed rules", func(t *testing.T) {
        g := cf.NewGrammar([]cf.Rule{
            {
                Left:  cf.Start,
                Right: "aSbS",
            },
            {
                Left:  cf.Start,
                Right: "",
            },
        })
        NewChomskyNormalizer(g).removeMixed()
        want := []cf.Rule{
            {
                Left:  'A',
                Right: "a",
            },
            {
                Left:  'B',
                Right: "b",
            },
            {
                Left:  cf.Start,
                Right: "ASBS",
            },
            {
                Left:  cf.Start,
                Right: "",
            },
        }
        assert.Equal(t, want, g.Rules)
    })
}

func TestChomskyNormalizer_removeLong(t *testing.T) {
    t.Run("Long rules", func(t *testing.T) {
        g := cf.NewGrammar([]cf.Rule{
            {
                Left:  cf.Start,
                Right: "WXYZ",
            },
            {
                Left:  'W',
                Right: "w",
            },
            {
                Left:  'X',
                Right: "x",
            },
            {
                Left:  'Y',
                Right: "y",
            },
            {
                Left:  'Z',
                Right: "z",
            },
            {
                Left:  cf.Start,
                Right: "",
            },
        })
        NewChomskyNormalizer(g).removeLong()
        want := []cf.Rule{
            {
                Left:  cf.Start,
                Right: "WB",
            },
            {
                Left:  'A',
                Right: "YZ",
            },
            {
                Left:  'B',
                Right: "XA",
            },
            {
                Left:  'W',
                Right: "w",
            },
            {
                Left:  'X',
                Right: "x",
            },
            {
                Left:  'Y',
                Right: "y",
            },
            {
                Left:  'Z',
                Right: "z",
            },
            {
                Left:  cf.Start,
                Right: "",
            },
        }
        checkSet := tools.NewSet[cf.Rule](want...)
        for _, rule := range g.Rules {
            assert.True(t, checkSet.Has(rule), rule)
            checkSet.Delete(rule)
        }
        assert.True(t, checkSet.Empty())
    })
}

func TestChomskyNormalizer_removeNullProductive(t *testing.T) {
    t.Run("Null-productive rules", func(t *testing.T) {
        g := cf.NewGrammar([]cf.Rule{
            {
                Left:  cf.Start,
                Right: "AB",
            },
            {
                Left:  'A',
                Right: "",
            },
            {
                Left:  'A',
                Right: "a",
            },
            {
                Left:  'B',
                Right: "b",
            },
        })
        NewChomskyNormalizer(g).removeNullProductive()
        want := []cf.Rule{
            {
                Left:  cf.Start,
                Right: "AB",
            },
            {
                Left:  cf.Start,
                Right: "B",
            },
            {
                Left:  'A',
                Right: "a",
            },
            {
                Left:  'B',
                Right: "b",
            },
        }
        checkSet := tools.NewSet[cf.Rule](want...)
        for _, rule := range g.Rules {
            assert.True(t, checkSet.Has(rule), rule)
            checkSet.Delete(rule)
        }
        assert.True(t, checkSet.Empty())
    })
}

func TestChomskyNormalizer_removeUnit(t *testing.T) {
    t.Run("Unit rules", func(t *testing.T) {
        g := cf.NewGrammar([]cf.Rule{
            {
                Left:  cf.Start,
                Right: "A",
            },
            {
                Left:  'A',
                Right: "B",
            },
            {
                Left:  'B',
                Right: "a",
            },
        })
        NewChomskyNormalizer(g).removeUnit()
        want := []cf.Rule{
            {
                Left:  cf.Start,
                Right: "a",
            },
        }
        assert.Equal(t, want, g.Rules)
    })
}

func TestChomskyNormalizer_isNormalForm(t *testing.T) {
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
            want: true,
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
            assert.Equal(t, tc.want, NewChomskyNormalizer(g).isNormalForm())
        })
    }
}
