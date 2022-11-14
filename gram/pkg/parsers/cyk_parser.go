package parsers

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/gram/internal/common"
    "mipt_formal/gram/internal/modify"

    impl "mipt_formal/gram/internal/parsers"
)

func NewCYKParser(rules []string) (Parser, error) {
    grammar, err := cf.ParseGrammar(common.NewStringsReader(rules))
    if err != nil {
        return nil, err
    }
    modify.NewChomskyNormalizer(grammar).Run()
    return &cykParser{
        grammar: grammar,
    }, nil
}

type cykParser struct {
    grammar *cf.Grammar
}

func (c *cykParser) Check(word string) bool {
    return impl.CYK(c.grammar, word)
}
