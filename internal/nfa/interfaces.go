package nfa

type Compiler interface {
    Compile(expr string) (*Machine, error)
}
