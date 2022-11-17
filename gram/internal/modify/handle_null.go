package modify

import "mipt_formal/gram/internal/cf"

func (n *ChomskyNormalizer) handleNull() {
    if !n.needHandleNull {
        return
    }

    n.grammar.Rules = append(n.grammar.Rules, cf.Rule{
        Left:  cf.Start,
        Right: string([]byte{cf.Epsilon}),
    })
}
