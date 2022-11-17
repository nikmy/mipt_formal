package modify

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

func (n *ChomskyNormalizer) removeUnit() {
    children := make(map[byte][][]byte, len(n.grammar.Rules))
    for _, rule := range n.grammar.Rules {
        children[rule.Left] = append(children[rule.Left], []byte(rule.Right))
    }

    newRules := make([]cf.Rule, 0, len(n.grammar.Rules))
    for left, rights := range children {
        newRights := make([][]byte, 0, len(rights))
        for _, right := range rights {
            if len(right) == 1 && !cf.IsTerminal(right[0]) {
                q := tools.NewQueue[byte](right[0])
                for !q.Empty() {
                    next := q.Pop()
                    for _, transitiveRight := range children[next] {
                        if len(transitiveRight) > 1 {
                            newRights = append(newRights, transitiveRight)
                            continue
                        }
                        if cf.IsTerminal(transitiveRight[0]) {
                            newRights = append(newRights, transitiveRight)
                            continue
                        }
                        if transitiveRight[0] == next { // avoid loops
                            continue
                        }
                        q.Push(transitiveRight[0])
                    }
                }
            } else {
                newRights = append(newRights, right)
            }
        }
        for _, newRight := range newRights {
            newRules = append(newRules, cf.Rule{
                Left:  left,
                Right: string(newRight),
            })
        }
    }

    n.grammar.Rules = newRules

    n.removeUnreachable()
}
