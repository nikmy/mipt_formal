package modify

import "mipt_formal/gram/internal/cf"

func (n *ChomskyNormalizer) addProxyStart() {
    startAlias := n.nonTerminalsFreeList.Pop()
    newRules := make([]cf.Rule, 0, len(n.grammar.Rules)+1)
    newRules = append(newRules, cf.Rule{
        Left:  cf.Start,
        Right: string(startAlias),
    })
    for _, rule := range n.grammar.Rules {
        var newLeft byte
        newRight := make([]byte, 0, len(rule.Right))
        if rule.Left == cf.Start {
            newLeft = startAlias
        }
        for _, symbol := range []byte(rule.Right) {
            if symbol == cf.Start {
                newRight = append(newRight, startAlias)
            } else {
                newRight = append(newRight, symbol)
            }
        }
        newRules = append(newRules, cf.Rule{
            Left:  newLeft,
            Right: string(newRight),
        })
    }
    n.grammar.Rules = newRules
    n.startAlias = startAlias
}
