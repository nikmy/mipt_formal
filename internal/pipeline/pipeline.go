package pipeline

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/doa"
    "mipt_formal/internal/modify"
    "mipt_formal/internal/nfa"
)

func New(logger common.Logger, r common.Reader, c nfa.Compiler, steps []modify.Step, w doa.Writer) func() {
    return func() {
        input, err := r.Read()
        if err != nil {
            logger.Fatal(err)
        }
        machine, err := c.Compile(input)
        if err != nil {
            logger.Fatal(err)
        }

        modify.Sequence(logger, steps...)(machine)

        err = w.Write(machine)
        if err != nil {
            logger.Fatal(err)
        }
    }
}
