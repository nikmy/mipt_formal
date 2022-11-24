package parsers

import (
    "github.com/stretchr/testify/assert"
    "mipt_formal/gram/internal/cf"
    "testing"
)

func TestEarley(t *testing.T) {
    type args struct {
        rules []cf.Rule
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
                rules: []cf.Rule{
                    {
                        Left:  cf.Start,
                        Right: "",
                    },
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
            name: "correct bracket sequence",
            args: args{
                rules: []cf.Rule{
                    {
                        Left:  cf.Start,
                        Right: "aSbS",
                    },
                    {
                        Left:  cf.Start,
                        Right: "",
                    },
                },
                words: []string{
                    "aabb",
                    "abba",
                    "ababaabb",
                    "ba",
                    "aaabbabb",
                    "aaaaaaabbbbbbb",
                    "aaabaaa",
                },
            },
            want: []bool{
                true,
                false,
                true,
                false,
                true,
                true,
                false,
            },
        },
        {
            name: "just works",
            args: args{
                rules: []cf.Rule{
                    {
                        Left:  cf.Start,
                        Right: "AT",
                    },
                    {
                        Left:  'A',
                        Right: "aA",
                    },
                    {
                        Left:  'A',
                        Right: "",
                    },
                    {
                        Left:  'T',
                        Right: "aTb",
                    },
                    {
                        Left:  'T',
                        Right: "",
                    },
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
                rules: []cf.Rule{
                    {
                        Left:  cf.Start,
                        Right: "abcd",
                    },
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
            g := cf.NewGrammar(tc.args.rules)
            got := make([]bool, 0, len(tc.args.words))
            for _, word := range tc.words {
                got = append(got, Earley(g, word))
            }
            for i := range got {
                assert.Equal(t, tc.want[i], got[i], tc.words[i])
            }
        })
    }
}
