package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func (n *ChomskyNormalizer) removeMixed() {
    producers := make(map[byte]byte) // Sigma -> N
    N := tools.NewStack[byte](n.nonTerminalsFreeList.AsSlice()...)

    for i, rule := range n.grammar.Rules {
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
                        n.nonTerminalsFreeList.Delete(producer)
                    }
                    newRight = append(newRight, producer)
                }
            }
            n.grammar.Rules[i].Right = string(newRight)
        }
    }
}
