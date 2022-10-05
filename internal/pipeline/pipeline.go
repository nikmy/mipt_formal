package pipeline

import (
	"mipt_formal/internal/common"
	"mipt_formal/internal/doa"
	"mipt_formal/internal/modify"
	"mipt_formal/internal/nfa"
)

func New(logger common.Logger, r common.Reader, c nfa.Compiler, steps []modify.Step, w doa.Writer) func() {
	return func() {
		logger.Info("Reading input...")
		input, err := r.Read()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Compiling NFA...")
		machine, err := c.Compile(input)
		if err != nil {
			logger.Fatal(err)
		}

		modify.Sequence(logger, steps...)(machine)

		logger.Info("Writing input...")
		err = w.Write(machine)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Done")
	}
}
