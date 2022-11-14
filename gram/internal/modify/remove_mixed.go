package modify

import (
    "mipt_formal/gram/internal/cf"
)

func (n *ChomskyNormalizer) removeMixed() {
    producers := make(map[byte]byte) // Sigma -> N
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) == 1 && cf.IsTerminal(rule.Right[0]) {
            producers[rule.Right[0]] = rule.Left
        }
    }
    N := n.nonTerminalsFreeList

    newRules := make([]cf.Rule, 0, len(n.grammar.Rules))
    for _, rule := range n.grammar.Rules {
        hasTerminals, hasNonTerminals, isMixed := false, false, false
        for _, symbol := range []byte(rule.Right) {
            if cf.IsTerminal(symbol) {
                hasTerminals = true
                if hasNonTerminals {
                    isMixed = true
                    break
                }
            } else {
                hasNonTerminals = true
                if hasTerminals {
                    isMixed = true
                    break
                }
            }
        }

        if !isMixed {
            newRules = append(newRules, rule)
            continue
        }

        newRight := make([]byte, 0, len(rule.Right))
        for _, symbol := range []byte(rule.Right) {
            if !cf.IsTerminal(symbol) {
                newRight = append(newRight, symbol)
                continue
            }

            producer, known := producers[symbol]
            if !known {
                if N.Empty() {
                    panic("not enough symbols is non-terminals alphabet")
                }
                producer = N.Pop()
                producers[symbol] = producer
                newRules = append(newRules, cf.Rule{
                    Left:  producer,
                    Right: string(symbol),
                })
            }
            newRight = append(newRight, producer)
        }

        newRules = append(newRules, cf.Rule{
            Left:  rule.Left,
            Right: string(newRight),
        })
    }
    n.grammar.Rules = newRules
}
