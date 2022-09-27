package regex

import (
    "errors"
    "fmt"
    "strings"

    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
    "mipt_formal/internal/tools"
)

func NewCompiler() nfa.Compiler {
    return compiler{}
}

type compiler struct{}

func (c compiler) Compile(expr string) (*nfa.Machine, error) {
    re, err := c.postfix(expr)
    if err != nil {
        return nil, fmt.Errorf("regexp parsing error: %w", err)
    }

    fragments := tools.NewStack[fragment]()

    for _, sym := range []byte(re) {
        var start, accept *intrusiveState
        switch sym {
        case orOperator:
            if fragments.Size() < 2 {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, orOperator))
            }

            f1 := fragments.Pop()
            f2 := fragments.Pop()

            start = NewIntrusiveState(common.Epsilon, f1.Start, f2.Start)
            accept = NewIntrusiveState(common.Epsilon)

            f1.Accept.precede(accept)
            f2.Accept.precede(accept)
        case kleeneStar:
            if fragments.Empty() {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, kleeneStar))
            }

            f := fragments.Pop()
            start = NewIntrusiveState(common.Epsilon, f.Start, f.Accept)
            f.Accept.precede(start)
            accept = f.Accept
        case oneOrMore:
            if fragments.Empty() {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, oneOrMore))
            }

            f := fragments.Pop()
            f.Accept.precede(f.Start)
            start, accept = f.Start, f.Accept
        default:
            // hack: accept is always epsilon-labeled
            accept = NewIntrusiveState(common.Epsilon)
            start = NewIntrusiveState(common.Word(sym), accept)
        }

        fragments.Push(fragment{
            Start:  start,
            Accept: accept,
        })
    }

    for fragments.Size() > 1 {
        f := fragments.Pop()
        p := fragments.Pop()
        p.Accept.precede(f.Start)
        fragments.Push(fragment{
            Start:  p.Start,
            Accept: f.Accept,
        })
    }

    if fragments.Empty() {
        return nil, nil
    }

    return c.fragmentToMachine(fragments.Pop()), nil
}

func (compiler) fragmentToMachine(f fragment) *nfa.Machine {
    return nfa.New(RunDFSWalker(f.Start, f.Accept))
}

func (compiler) postfix(infix string) (string, error) {
    var result strings.Builder

    ops := tools.NewStack[byte]()
    for _, cur := range []byte(infix) {
        switch cur {
        case lBracket:
            ops.Push(cur)
        case rBracket:
            for {
                if ops.Empty() {
                    return "", errors.New("invalid parentheses")
                }
                op := ops.Pop()
                if op == lBracket {
                    break
                }
                result.WriteByte(op)
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
                if priority[cur] >= priority[ops.Top()] {
                    break
                }
                result.WriteByte(ops.Pop())
            }
            ops.Push(cur)
        default:
            result.WriteByte(cur)
        }
    }

    for !ops.Empty() {
        result.WriteByte(ops.Pop())
    }

    return result.String(), nil
}
