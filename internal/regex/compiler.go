package regex

import (
    "errors"
    "fmt"
    "strings"

    "mipt_formal/internal/fsm"
    "mipt_formal/internal/tools"
)

func NewCompiler() Compiler {
    return compiler{}
}

type compiler struct{}

func (c compiler) Compile(expr string) (fsm.Machine, error) {
    re, err := c.postfix(expr)
    if err != nil {
        return nil, fmt.Errorf("regexp parsing error: %w", err)
    }

    fragments := tools.NewStack[fragment]()

    for _, sym := range []byte(re) {
        var start, accept *IntrusiveState
        switch sym {
        case orOperator:
            if fragments.Size() < 2 {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, orOperator))
            }

            f1 := fragments.Pop()
            f2 := fragments.Pop()

            accept = NewIntrusiveState(fsm.Epsilon)
            start = NewIntrusiveState(fsm.Epsilon, f1.Start, f2.Start)

            f1.Accept.precede(accept)
            f2.Accept.precede(accept)
        case kleeneStar:
            if fragments.Empty() {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, kleeneStar))
            }

            f := fragments.Pop()
            accept = NewIntrusiveState(fsm.Epsilon)
            start = NewIntrusiveState(fsm.Epsilon, f.Start, accept)
            f.Accept.precede(f.Start, accept)
        case oneOrMore:
            if fragments.Empty() {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, oneOrMore))
            }

            f := fragments.Pop()
            accept = NewIntrusiveState(fsm.Epsilon)
            start = NewIntrusiveState(fsm.Epsilon, f.Start)

            f.Accept.precede(f.Start, accept)
        default:
            accept = NewIntrusiveState(fsm.Epsilon)
            start = NewIntrusiveState(fsm.Word(sym), accept)
        }

        fragments.Push(fragment{
            Start:  start,
            Accept: accept,
        })
    }

    if fragments.Size() != 1 {
        return nil, errors.New("too many fragments")
    }

    return c.fragmentToMachine(fragments.Pop()), nil
}

func (compiler) fragmentToMachine(f fragment) fsm.Machine {
    return fsm.NewNFA(RunDFSWalker(f.Start, f.Accept))
}

func (compiler) postfix(infix string) (string, error) {
    var result strings.Builder

    ops := tools.NewStack[byte]()
    for _, cur := range []byte(infix) {
        switch cur {
        case lBracket:
            ops.Push(cur)
        case rBracket:
            for ops.Pop() != lBracket {
                if ops.Empty() {
                    return "", errors.New("invalid parentheses")
                }
            }
        case kleeneStar:
            fallthrough
        case orOperator:
            fallthrough
        case oneOrMore:
            for {
                if ops.Empty() {
                    break
                }
                last := ops.Pop()
                if priority[cur] >= priority[last] {
                    break
                }
                result.WriteByte(last)
            }
            fallthrough
        default:
            result.WriteByte(cur)
        }
    }

    for !ops.Empty() {
        result.WriteByte(ops.Pop())
    }

    return result.String(), nil
}
