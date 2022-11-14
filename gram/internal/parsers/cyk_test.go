package parsers

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/gram/internal/common"
    "testing"
)

func TestCYK(t *testing.T) {
    type args struct {
        rules []string
        words []string
    }
    type testcase struct {
        name string
        args
        want []bool
    }

    cases := [...]testcase{
        {
            name: "empty word",
            args: args{
                rules: []string{
                    // S -> Epsilon
                    "S -> _",
                },
                words: []string{
                    "a",
                    "",
                },
            },
            want: []bool{
                false,
                true,
            },
        },
        {
            // S -> aSbS | Epsilon
            name: "correct bracket sequence",
            args: args{
                rules: []string{
                    "S -> AR",
                    "R -> ST",
                    "T -> BS",
                    "A -> a",
                    "B -> b",
                    "S -> _",
                },
                words: []string{
                    "aabb",
                    "abba",
                    "ababaabb",
                    "ba",
                },
            },
            want: []bool{
                true,
                false,
                true,
                false,
            },
        },
        {
            // S -> AT
            // A -> aA  | Epsilon
            // T -> aTb | Epsilon
            name: "just works",
            args: args{
                rules: []string{
                    "S -> TR",
                    "T -> AT",
                    "T -> _",
                    "A -> a",
                    "R -> AU",
                    "R -> _",
                    "U -> RB",
                    "B -> b",
                    "S -> _",
                },
                words: []string{
                    "aaaaabb",
                    "baba",
                    "aaabb",
                    "ba",
                },
            },
            want: []bool{
                true,
                false,
                true,
                false,
            },
        },
        {
            name: "just works 2",
            args: args{
                rules: []string{
                    // S -> abcd
                    "S -> AT",
                    "U -> CD",
                    "A -> a",
                    "T -> BU",
                    "D -> d",
                    "B -> b",
                    "C -> c",
                },
                words: []string{
                    "abcd",
                    "abba",
                    "aaabb",
                    "a",
                    "",
                },
            },
            want: []bool{
                true,
                false,
                false,
                false,
                false,
            },
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            g, _ := cf.ParseGrammar(common.NewStringsReader(tc.args.rules))
            got := make([]bool, 0, len(tc.args.words))
            for _, word := range tc.words {
                got = append(got, CYK(g, word))
            }
        })
    }
}
