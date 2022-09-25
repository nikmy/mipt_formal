package regex

import "mipt_formal/internal/fsm"

type Compiler interface {
    Compile(expr string) (fsm.Machine, error)
}
