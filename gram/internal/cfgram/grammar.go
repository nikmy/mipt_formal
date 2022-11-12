package cfgram

import (
    "mipt_formal/gram/internal/common"
)

func NewGrammar(rules []Rule) *Grammar {
    return &Grammar{
        Rules: rules,
    }
}

type Grammar struct {
    Rules []Rule
}

// NonTerminal is type for non-terminals. It is separated for static check for context-freedom.
type NonTerminal byte

// Symbol is common type for terminals and non-terminals
type Symbol byte

func (s Symbol) NonTerminal() (NonTerminal, error) {
    if !isNonTerminal(byte(s)) {
        return 0, common.IsNotNonTerminalSymbolError
    }
    return NonTerminal(s), nil
}

type Rule struct {
    Left  NonTerminal
    Right []Symbol
}
