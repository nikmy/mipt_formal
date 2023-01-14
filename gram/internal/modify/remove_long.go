package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func (n *ChomskyNormalizer) removeLong() {
    producers := make(map[string]byte) // Sigma x Sigma -> N
    N := n.nonTerminalsFreeList

    newRules := make([]cf.Rule, 0, len(n.grammar.Rules))
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) <= 2 {
            newRules = append(newRules, rule)
            continue
        }

        // fill rules chain from the end to avoid repeats
        toSplit := tools.NewStack[byte]()
        for i := 0; i < len(rule.Right)-2; i++ {
            toSplit.Push(rule.Right[i])
        }

        nextRight := rule.Right[len(rule.Right)-2:]
        for {
            producer, known := producers[nextRight]
            if toSplit.Empty() {
                producer, known = rule.Left, true
            }
            if !known {
                if N.Empty() {
                    panic("not enough symbols is non-terminals alphabet")
                }
                producer = N.Pop()
                producers[nextRight] = producer
            }
            newRules = append(newRules, cf.Rule{
                Left:  producer,
                Right: nextRight,
            })

            if toSplit.Empty() {
                break
            }

            nextRight = string([]byte{toSplit.Pop(), producer})
        }
    }
    n.grammar.Rules = newRules
}
