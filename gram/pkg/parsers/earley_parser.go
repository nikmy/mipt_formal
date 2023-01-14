package parsers

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/gram/internal/common"
    impl "mipt_formal/gram/internal/parsers"
)

func NewEarleyParser(rules []string) (Parser, error) {
    grammar, err := cf.ParseGrammar(common.NewStringsReader(rules))
    if err != nil {
        return nil, err
    }
    return &earleyParser{
        grammar: grammar,
    }, nil
}

type earleyParser struct {
    grammar *cf.Grammar
}

func (e *earleyParser) Check(word string) bool {
    return impl.Earley(e.grammar, word)
}
