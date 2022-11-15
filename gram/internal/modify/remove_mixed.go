package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func removeMixed(g *cf.Grammar) {
    availableNonTerminals := tools.NewSet[byte]([]byte("ABCDEFGHIJKLMNOPQRTUVWXYZ")...)
    for _, rule := range g.Rules {
        if availableNonTerminals.Has(rule.Left) {
            availableNonTerminals.Delete(rule.Left)
        }
    }
    N := tools.StackFromSlice(availableNonTerminals.AsSlice())
    availableNonTerminals = nil
    producers := make(map[byte]byte) // Sigma -> N

    for _, rule := range g.Rules {
        hasTerminals, hasNonTerminals, isMixed := false, false, false
        for _, symbol := range []byte(rule.Right) {
            if cf.IsNonTerminal(symbol) {
                hasNonTerminals = true
                if hasTerminals {
                    isMixed = true
                    break
                }
            } else {
                hasTerminals = true
                if hasNonTerminals {
                    isMixed = true
                    break
                }
            }
        }

        if isMixed {
            newRight := make([]byte, 0, len(rule.Right))
            for _, symbol := range []byte(rule.Right) {
                if !cf.IsNonTerminal(symbol) {
                    producer, known := producers[symbol]
                    if !known {
                        if N.Empty() {
                            panic("not enough symbols is non-terminals alphabet")
                        }
                        producer = N.Pop()
                        producers[symbol] = producer
                    }
                    newRight = append(newRight, producer)
                }
            }
            rule.Right = string(newRight)
        }
    }
}
