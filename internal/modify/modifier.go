package modify

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/nfa"
)

type Modifier func(*nfa.Machine)

type Step struct {
    Name string
    Func Modifier
}

func Sequence(logger common.Logger, steps ...Step) Modifier {
    return func(m *nfa.Machine) {
        if m == nil {
            logger.Info("Empty input")
            return
        }
        logger.Info("Running NFA modifying sequence...")
        for i, step := range steps {
            logger.Info(m.DOA())
            logger.Infof("Step %d: %s...", i+1, step.Name)
            step.Func(m)
        }
        logger.Info("Done")
    }
}
