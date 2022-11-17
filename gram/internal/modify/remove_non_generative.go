package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func (n *ChomskyNormalizer) removeNonGenerative() {
    generative := tools.NewSet[byte]()
    nonGenRight := make([]tools.Set[byte], 0, len(n.grammar.Rules))

    for i, rule := range n.grammar.Rules {
        nonGenRight = append(nonGenRight, tools.NewSet[byte]())
        for _, symbol := range []byte(rule.Right) {
            if !cf.IsTerminal(symbol) && !generative.Has(symbol) {
                nonGenRight[i].Insert(symbol)
            }
        }
        if nonGenRight[i].Empty() {
            generative.Insert(rule.Left)
        }
    }

    newGenerative := !generative.Empty()
    for newGenerative {
        newGenerative = false
        for i, rule := range n.grammar.Rules {
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
            if nonGenRight[i].Empty() {
                generative.Insert(rule.Left)
                newGenerative = true
            }
        }
    }

    newRules := make([]cf.Rule, 0, len(n.grammar.Rules))
    for _, rule := range n.grammar.Rules {
        if !generative.Has(rule.Left) {
            continue
        }
        skipRule := false
        for _, symbol := range []byte(rule.Right) {
            if !cf.IsTerminal(symbol) && !generative.Has(symbol) {
                skipRule = true
                break
            }
        }

        if !skipRule {
            newRules = append(newRules, rule)
        }
    }
    n.grammar.Rules = newRules
}
