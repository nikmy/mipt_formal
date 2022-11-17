package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func NewChomskyNormalizer(g *cf.Grammar) *ChomskyNormalizer {
    needHandleNull := false
    freelist := tools.NewSet[byte]([]byte("ABCDEFGHIJKLMNOPQRTUVWXYZ")...)
    for _, rule := range g.Rules {
        if rule.Left == cf.Start && len(rule.Right) == 1 && rule.Right[0] == cf.Epsilon {
            needHandleNull = true
            break
        }
        if freelist.Has(rule.Left) {
            freelist.Delete(rule.Left)
        }
    }

    return &ChomskyNormalizer{
        grammar:              g,
        needHandleNull:       needHandleNull,
        nonTerminalsFreeList: freelist,
    }
}

type ChomskyNormalizer struct {
    grammar              *cf.Grammar
    needHandleNull       bool
    nonTerminalsFreeList tools.Set[byte]
}

func (n *ChomskyNormalizer) Run() {
    if n == nil || n.grammar == nil || n.checkNF() { // lazy
        return
    }
    n.removeNonGenerative()
    n.removeUnreachable()
    n.removeMixed()
    n.removeLong()
    n.removeNullProductive()
    n.handleNull()
    n.removeUnit()
}

func (n *ChomskyNormalizer) checkNF() bool {
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) > 2 {
            return false
        }
        if len(rule.Right) == 1 {
            if cf.IsNonTerminal(rule.Right[0]) {
                return false
            }
        } else {
            if !cf.IsNonTerminal(rule.Right[0]) || !cf.IsNonTerminal(rule.Right[1]) {
                return false
            }
        }
    }
    return true
}
