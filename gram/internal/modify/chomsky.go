package modify

import (
    "sort"

    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func NewChomskyNormalizer(g *cf.Grammar) *ChomskyNormalizer {
    needHandleNull := false
    freeNonTerminals := tools.NewSet[byte]([]byte(cf.NonTerminalsAlphabet)...)
    for _, rule := range g.Rules {
        if rule.Left == cf.Start && len(rule.Right) == 0 {
            needHandleNull = true
            break
        }
        if freeNonTerminals.Has(rule.Left) {
            freeNonTerminals.Delete(rule.Left)
        }
    }

    freelist := freeNonTerminals.AsSlice()
    sort.Slice(freelist, func(i, j int) bool {
        return freelist[i] < freelist[j]
    })

    return &ChomskyNormalizer{
        grammar:              g,
        needHandleNull:       needHandleNull,
        nonTerminalsFreeList: tools.NewQueue[byte](freelist...),
    }
}

type ChomskyNormalizer struct {
    grammar              *cf.Grammar
    needHandleNull       bool
    nonTerminalsFreeList *tools.Queue[byte]
    startAlias           byte
}

func (n *ChomskyNormalizer) Run() {
    if n == nil || n.grammar == nil || n.isNormalForm() { // lazy
        return
    }
    n.addProxyStart()
    n.removeNonGenerative()
    n.removeUnreachable()
    n.removeMixed()
    n.removeLong()
    n.removeNullProductive()
    n.handleNull()
    n.removeUnit()

    n.grammar.Rules = n.grammar.Rules[:len(n.grammar.Rules):len(n.grammar.Rules)] // shrink to fit
}

func (n *ChomskyNormalizer) isNormalForm() bool {
    for _, rule := range n.grammar.Rules {
        if len(rule.Right) == 0 {
            if rule.Left != cf.Start {
                return false
            }
            continue
        }

        if len(rule.Right) == 1 && cf.IsTerminal(rule.Right[0]) {
            continue
        }
        if len(rule.Right) == 2 && !cf.IsTerminal(rule.Right[0]) && !cf.IsTerminal(rule.Right[1]) {
            continue
        }
        return false
    }
    return true
}
