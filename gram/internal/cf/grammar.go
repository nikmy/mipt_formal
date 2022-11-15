package cf

func NewGrammar(rules []Rule) *Grammar {
    return &Grammar{
        Rules: rules,
    }
}

type Grammar struct {
    Rules []Rule
}

type Rule struct {
    Left  byte
    Right string
}

func IsNonTerminal(symbol byte) bool {
    return symbol >= 'A' && symbol <= 'Z'
}
