package regex

import "mipt_formal/internal/nfa"

type Compiler interface {
    Compile(expr string) (*nfa.Machine, error)
}
