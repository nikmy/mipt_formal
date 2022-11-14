package parsers

import (
    "mipt_formal/gram/internal/cf"
)

func CYK(g *cf.Grammar, word string) bool {
    dp := make(map[byte][][]bool, len(g.Rules))

    singles := make(map[byte]byte)
    bins := make([]cf.Rule, 0, 1)

    for _, rule := range g.Rules {
        if _, ok := dp[rule.Left]; !ok {
            dp[rule.Left] = make([][]bool, len(word))
            for i := range word {
                dp[rule.Left][i] = make([]bool, len(word)+1)
            }
        }
        if len(rule.Right) == 1 {
            if rule.Right[0] != cf.Epsilon {
                singles[rule.Right[0]] = rule.Left
            } else if len(word) == 0 {
                // For 'S -> Epsilon' rule
                return true
            }
        } else {
            bins = append(bins, rule)
        }
    }

    if len(word) == 0 {
        return false
    }

    // fill base cases
    for i, s := range []byte(word) {
        if producer, ok := singles[s]; ok {
            dp[producer][i][i] = true
        }
    }

    for m := 1; m < len(word); m++ {
        for i := 0; i < len(word)-m; i++ {
            j := i + m

        loop:
            for _, binRule := range bins {
                A := binRule.Left
                B := binRule.Right[0]
                C := binRule.Right[1]
                for k := i; k < j; k++ {
                    if dp[B][i][k] && dp[C][k+1][j] {
                        dp[A][i][j] = true
                        break loop
                    }
                }
            }
        }
    }

    return dp[cf.Start][0][len(word)]
}
