package modify

import "mipt_formal/gram/internal/cf"

func NewChomskyNormalizer(g *cf.Grammar) *ChomskyNormalizer {
    needHandleNull := false
    for _, rule := range g.Rules {
        if rule.Left == cf.Start && len(rule.Right) == 1 && rule.Right[0] == cf.Epsilon {
            needHandleNull = true
            break
        }
    }

    return &ChomskyNormalizer{
        grammar: g,
        modifiers: []func(*cf.Grammar){
            removeNonGenerative,
            removeUnreachable,
            removeMixed,
            removeLong,
            removeNullProductive,
            handleNull,
            removeUnit,
        },
        handleNull: needHandleNull,
    }
}

type ChomskyNormalizer struct {
    grammar    *cf.Grammar
    modifiers  []func(*cf.Grammar)
    handleNull bool
}

func (n *ChomskyNormalizer) Run() {
    for i, step := range n.modifiers {
        if i == len(n.modifiers)-2 && !n.handleNull {
            continue
        }
        step(n.grammar)
    }
}
