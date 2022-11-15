package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func removeNonGenerative(g *cf.Grammar) {
    generative := tools.NewSet[byte]()
    nonGenRight := make([]tools.Set[byte], 0, len(g.Rules))

    for i, rule := range g.Rules {
        nonGenRight = append(nonGenRight, tools.NewSet[byte]())
        for _, symbol := range []byte(rule.Right) {
            if cf.IsNonTerminal(symbol) && !generative.Has(symbol) {
                nonGenRight[i].Insert(symbol)
            }
        }
        if nonGenRight[i].IsEmpty() {
            generative.Insert(rule.Left)
        }
    }

    newGenerative := !generative.IsEmpty()
    for newGenerative {
        for i, rule := range g.Rules {
            if generative.Has(rule.Left) {
                continue
            }

            toDelete := make([]byte, 0)
            for symbol := range nonGenRight[i] {
                if generative.Has(symbol) {
                    toDelete = append(toDelete, symbol)
                }
            }
            for _, symbol := range toDelete {
                nonGenRight[i].Delete(symbol)
            }
            if nonGenRight[i].IsEmpty() {
                generative.Insert(rule.Left)
                newGenerative = true
            }
        }
    }
}
