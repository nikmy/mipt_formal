package parsers

import (
    "mipt_formal/gram/internal/cf"
)

func CYK(g *cf.Grammar, word string) bool {
    reach := make(map[byte][][]bool, len(g.Rules))

    singles := make(map[byte][]byte)
    bins := make([]cf.Rule, 0, 1)

    for _, rule := range g.Rules {
        if _, ok := reach[rule.Left]; !ok {
            reach[rule.Left] = make([][]bool, len(word))
            for i := range word {
                reach[rule.Left][i] = make([]bool, len(word))
            }
        }

        switch len(rule.Right) {
        case 2:
            bins = append(bins, rule)
        case 1:
            singles[rule.Right[0]] = append(singles[rule.Right[0]], rule.Left)
        case 0:
            if len(word) == 0 {
                return true
            }
        }
    }

    if len(word) == 0 {
        return false
    }

    // fill base cases
    for i, s := range []byte(word) {
        if producers, ok := singles[s]; ok {
            for _, producer := range producers {
                reach[producer][i][i] = true
            }
        }
    }

    for substringLen := 1; substringLen < len(word); substringLen++ {
        for i := 0; i < len(word)-substringLen; i++ {
            j := i + substringLen

            for _, binRule := range bins {
                A := binRule.Left
                B := binRule.Right[0]
                C := binRule.Right[1]
                for k := i; k < j; k++ {
                    if reach[B][i][k] && reach[C][k+1][j] {
                        reach[A][i][j] = true
                        break
                    }
                }
            }
        }
    }

    return reach[cf.Start][0][len(word)-1]
}
