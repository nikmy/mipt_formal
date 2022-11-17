package modify

import "mipt_formal/gram/internal/cf"

func (n *ChomskyNormalizer) removeNullProductive() {
    newRules := make([]cf.Rule, 0, len(n.grammar.Rules))
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) == 0 && rule.Right[0] == cf.Epsilon {
            continue
        }
        newRules = append(newRules, rule)
    }
    n.grammar.Rules = newRules
}
