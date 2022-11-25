package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func (n *ChomskyNormalizer) removeNullProductive() {
    // find all such E that E -> null
    nullProductive := tools.NewSet[byte]()
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) == 0 {
            nullProductive.Insert(rule.Left)
            continue
        }
    }

    n.nullClosure(nullProductive)

    if nullProductive.Has(cf.Start) {
        n.needHandleNull = true
    }

    // A -> EB and A -> BE, E -> null cases + remove null-productive
    newRules := make([]cf.Rule, 0, len(n.grammar.Rules))
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) == 2 {
            leftNull := nullProductive.Has(rule.Right[0])
            rightNull := nullProductive.Has(rule.Right[1])
            if !leftNull && !rightNull {
                newRules = append(newRules, rule)
                continue
            }
            if leftNull {
                newRules = append(newRules, cf.Rule{
                    Left:  rule.Left,
                    Right: rule.Right[1:],
                })
            }
            if rightNull {
                newRules = append(newRules, cf.Rule{
                    Left:  rule.Left,
                    Right: rule.Right[:1],
                })
            }
            newRules = append(newRules, rule)
            continue
        }

        if len(rule.Right) > 2 {
            panic("couldn't remove null productive: unexpected right rule: " + rule.Right)
        }

        if len(rule.Right) == 1 {
            newRules = append(newRules, rule)
        }
    }

    n.grammar.Rules = newRules
}

func (n *ChomskyNormalizer) nullClosure(nullProductive tools.Set[byte]) {
    newNullProductive := !nullProductive.Empty()
    for newNullProductive {
        newNullProductive = false
        for i, rule := range n.grammar.Rules {
            if nullProductive.Has(rule.Left) {
                continue
            }
            if len(rule.Right) == 1 && nullProductive.Has(rule.Right[0]) {
                nullProductive.Insert(rule.Left)
                newNullProductive = true
                n.grammar.Rules[i].Right = ""
                continue
            }
        }
    }
}
