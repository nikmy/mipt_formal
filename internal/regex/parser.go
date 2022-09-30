package regex

import (
    "errors"
    "fmt"
)

type AST struct {
    operator func(...*fragment) *fragment
    children []*AST
}

func (t *AST) compile() *fragment {
    if t == nil {
        return nil
    }
    args := make([]*fragment, 0, len(t.children))
    for _, child := range t.children {
        arg := child.compile()
        if arg != nil {
            args = append(args, arg)
        }
    }
    return t.operator(args...)
}

type parser struct{}

func (p parser) parseRegex(re string) (*AST, error) {
    return p.splitByOr(re)
}

func (p parser) splitByOr(expr string) (*AST, error) {
    split := make([]string, 0)
    balance := 0
    l, r := 0, 0
    for r < len(expr) {
        c := expr[r]
        r++
        switch c {
        case lBracket:
            balance++
        case rBracket:
            balance--
            if balance < 0 {
                return nil, errors.New(invalidParenthesesError)
            }
        case orOperator:
            if balance == 0 {
                split = append(split, expr[l:r-1])
                l = r
            }
        }

        if r == len(expr) {
            split = append(split, expr[l:r])
        }
    }

    if len(split) < 2 {
        return p.eval(expr)
    }

    nodes := make([]*AST, 0, len(split))
    for _, expr := range split {
        node, err := p.eval(expr)
        if err != nil {
            return nil, err
        }
        if node != nil {
            nodes = append(nodes, node)
        }
    }

    if len(nodes) == 0 {
        return nil, nil
    }

    if len(nodes) == 1 {
        return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, string(orOperator)))
    }

    return &AST{
        operator: _or,
        children: nodes,
    }, nil
}

func (p parser) eval(expr string) (*AST, error) {
    if len(expr) == 0 {
        return nil, nil
    }

    nodes := make([]*AST, 0, len(expr))
    i := 0
    for i < len(expr) {
        c := expr[i]
        i++

        if c == lBracket {
            if i == len(expr) {
                return nil, errors.New("invalid parentheses")
            }

            j := i
            for i < len(expr) {
                if expr[i] == rBracket {
                    break
                }
                i++
            }

            if i == len(expr) {
                return nil, errors.New("invalid parentheses")
            }

            node, err := p.splitByOr(expr[j:i])
            if err != nil {
                return nil, err
            }
            nodes = append(nodes, node)
            i++
            continue
        }

        if c == rBracket {
            return nil, errors.New(invalidParenthesesError)
        }

        if c == kleeneStar {
            if len(nodes) == 0 {
                return nil, errors.New(fmt.Sprintf(fewArgumentsErrorFormat, string(kleeneStar)))
            }
            last := nodes[len(nodes)-1]
            nodes[len(nodes)-1] = &AST{
                operator: _kleene,
                children: []*AST{last},
            }
            continue
        }

        // if you are there, you are not special character
        nodes = append(nodes, &AST{
            operator: _id(c),
            children: nil,
        })
    }

    if len(nodes) == 1 {
        return nodes[0], nil
    }

    return &AST{
        operator: _concat,
        children: nodes,
    }, nil
}
