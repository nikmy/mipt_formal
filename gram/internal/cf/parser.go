package cf

import (
    "bytes"
    "io"

    "mipt_formal/gram/internal/common"
)

func ParseGrammar(in reader) (*Grammar, error) {
    rules := make([]Rule, 0)

    lineNumber := 1

    wrongRule := func(err error) error {
        if err != nil {
            return common.Wrapf(err, "wrong rule on line %d", lineNumber)
        }
        return common.Errorf("wrong rule on line %d", lineNumber)
    }

    for {
        line, err := in.ReadLine()
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, err
        }

        if len(line) == 0 {
            continue
        }

        parts := bytes.Split(line, []byte("->"))
        if len(parts) != 2 {
            return nil, wrongRule(nil)
        }

        left, right := bytes.TrimSpace(parts[0]), bytes.TrimSpace(parts[1])
        if len(left) != 1 || !isNonTerminal(left[0]) {
            return nil, wrongRule(common.Errorf("left side must be single non-terminal"))
        }
        if lineNumber == 1 && left[0] != Start {
            return nil, wrongRule(common.Error("first rule must contain S in left side"))
        }
        l := left[0]

        if len(right) == 0 {
            return nil, wrongRule(common.Errorf("right side must be not empty"))
        }

        r := make([]byte, 0, len(right))
        for _, sym := range right {
            if sym == Epsilon {
                continue
            }
            if !IsTerminal(sym) && !isNonTerminal(sym) {
                return nil, wrongRule(common.Errorf("unexpected symbol '%v'", sym))
            }
            r = append(r, sym)
        }

        rules = append(rules, Rule{
            Left:  l,
            Right: string(r),
        })

        lineNumber++
    }

    return NewGrammar(rules), nil
}

func isNonTerminal(symbol byte) bool {
    return symbol >= 'A' || symbol <= 'Z'
}
