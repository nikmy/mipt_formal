package parsers

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
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

func Test_scan(t *testing.T) {
    type args struct {
        history *situationSet
        j       int
        w       string
    }
    type result struct {
        delta   *situationSet
        history *situationSet
    }
    type testcase struct {
        name string
        args args
        want result
    }

    rules := []cf.Rule{
        {
            Left:  fakeStart,
            Right: string(cf.Start),
        },
        {
            Left:  cf.Start,
            Right: "aB",
        },
    }

    cases := [...]testcase{
        {
            name: "just works",
            args: args{
                history: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {
                            cf.Start: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &rules[0],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                            'a': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &rules[1],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                        {},
                    },
                    size: 2,
                },
                j: 0,
                w: "ab",
            },
            want: result{
                delta: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {},
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                    },
                    size: 1,
                },
                history: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {
                            cf.Start: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &rules[0],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                            'a': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &rules[1],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                    },
                    size: 3,
                },
            },
        },
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := scan(tc.args.history, tc.args.j, tc.args.w)
            assert.Equal(t, *tc.want.delta, *got)
            assert.Equal(t, *tc.want.history, *tc.args.history)
        })
    }
}

func Test_predict(t *testing.T) {
    type args struct {
        delta   *situationSet
        history *situationSet
        j       int
        w       string
    }
    type result struct {
        newDelta *situationSet
        history  *situationSet
    }
    type testcase struct {
        name string
        args args
        want result
    }

    grammar := cf.NewGrammar([]cf.Rule{
        {
            Left:  fakeStart,
            Right: string(cf.Start),
        },
        {
            Left:  cf.Start,
            Right: "aB",
        },
        {
            Left:  'B',
            Right: "b",
        },
    })

    cases := [...]testcase{
        {
            name: "just works",
            args: args{
                delta: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {},
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                    },
                    size: 1,
                },
                history: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {
                            cf.Start: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[0],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                            'a': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                    },
                    size: 3,
                },
                j: 1,
                w: "ab",
            },
            want: result{
                newDelta: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {},
                        {
                            'b': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[2],
                                    RulePos: 0,
                                    WordPos: 1,
                                },
                            ),
                        },
                        {},
                    },
                    size: 1,
                },
                history: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {
                            cf.Start: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[0],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                            'a': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                            'b': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[2],
                                    RulePos: 0,
                                    WordPos: 1,
                                },
                            ),
                        },
                        {},
                    },
                    size: 4,
                },
            },
        },
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            newDelta := newSituationSet(len(tc.args.w) + 1)
            predict(tc.args.delta, tc.args.history, tc.args.j, grammar, newDelta)
            assert.Equal(t, *newDelta, *tc.want.newDelta)
            assert.Equal(t, *tc.args.history, *tc.want.history)
        })
    }
}

func Test_complete(t *testing.T) {
    type args struct {
        delta   *situationSet
        history *situationSet
        j       int
        w       string
    }
    type result struct {
        newDelta *situationSet
        history  *situationSet
    }
    type testcase struct {
        name string
        args args
        want result
    }

    grammar := cf.NewGrammar([]cf.Rule{
        {
            Left:  fakeStart,
            Right: string(cf.Start),
        },
        {
            Left:  cf.Start,
            Right: "aB",
        },
        {
            Left:  'B',
            Right: "",
        },
    })

    cases := [...]testcase{
        {
            name: "just works",
            args: args{
                delta: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {},
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                            endOfString: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[2],
                                    RulePos: 0,
                                    WordPos: 1,
                                },
                            ),
                        },
                        {},
                    },
                    size: 2,
                },
                history: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {
                            cf.Start: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[0],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                            'a': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                            endOfString: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[2],
                                    RulePos: 0,
                                    WordPos: 1,
                                },
                            ),
                        },
                        {},
                    },
                    size: 4,
                },
                j: 1,
                w: "ab",
            },
            want: result{
                newDelta: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {},
                        {
                            endOfString: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 2,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                    },
                    size: 1,
                },
                history: &situationSet{
                    data: []map[byte]tools.Set[earleySituation]{
                        {
                            cf.Start: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[0],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                            'a': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 0,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {
                            'B': tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 1,
                                    WordPos: 0,
                                },
                            ),
                            endOfString: tools.NewSet[earleySituation](
                                earleySituation{
                                    Rule:    &grammar.Rules[2],
                                    RulePos: 0,
                                    WordPos: 1,
                                },
                                earleySituation{
                                    Rule:    &grammar.Rules[1],
                                    RulePos: 2,
                                    WordPos: 0,
                                },
                            ),
                        },
                        {},
                    },
                    size: 5,
                },
            },
        },
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            newDelta := newSituationSet(len(tc.args.w) + 1)
            complete(tc.args.delta, tc.args.history, tc.args.j, newDelta)
            assert.Equal(t, *tc.want.newDelta, *newDelta)
            assert.Equal(t, *tc.want.history, *tc.args.history)
        })
    }
}
