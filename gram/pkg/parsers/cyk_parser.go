package parsers

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/gram/internal/modify"

    impl "mipt_formal/gram/internal/parsers"
)

func NewCYKParser(grammar *cf.Grammar) Parser {
    modify.NewChomskyNormalizer(grammar).Run()
    return &cykParser{
        grammar: grammar,
    }
}

type cykParser struct {
    grammar *cf.Grammar
}

func (c *cykParser) Check(word string) bool {
    return impl.CYK(c.grammar, word)
}
