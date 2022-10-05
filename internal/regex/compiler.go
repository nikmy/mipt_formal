package regex

import (
	"errors"
	"fmt"
	"mipt_formal/internal/nfa"
)

func NewCompiler() nfa.Compiler {
	return compiler{}
}

type compiler struct{}

func (c compiler) Compile(expr string) (*nfa.Machine, error) {
	var parser parser
	ast, err := parser.parseRegex(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("regex parsing error: %s", err))
	}
	f := ast.compile()
	return c.fragmentToMachine(f), nil
}

func (compiler) fragmentToMachine(f *fragment) *nfa.Machine {
	if f == nil {
		return nil
	}
	return nfa.NewMachine(RunDFSWalker(f.Start, f.Accept))
}
