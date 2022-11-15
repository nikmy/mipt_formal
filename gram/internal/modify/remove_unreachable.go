package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func findUnreachableDFS(prodGraph map[byte][]byte, current byte, visited tools.Set[byte]) {
    if visited.Has(current) {
        return
    }
    visited.Insert(current)
    for _, next := range prodGraph[current] {
        findUnreachableDFS(prodGraph, next, visited)
    }
}

func removeUnreachable(g *cf.Grammar) {
    prodGraph := make(map[byte][]byte)
    visited := tools.NewSet[byte]()
    for _, rule := range g.Rules {
        if _, known := prodGraph[rule.Left]; !known {
            prodGraph[rule.Left] = make([]byte, 0)
        }
        for _, symbol := range []byte(rule.Right) {
            if cf.IsNonTerminal(symbol) {
                prodGraph[rule.Left] = append(prodGraph[rule.Left], symbol)
            }
        }
    }
    findUnreachableDFS(prodGraph, cf.Start, visited)
    if visited.Size() == len(g.Rules) {
        return
    }

    newRules := make([]cf.Rule, 0, visited.Size())
    for _, rule := range g.Rules {
        if visited.Has(rule.Left) {
            newRules = append(newRules, rule)
        }
    }
    g.Rules = newRules
}
