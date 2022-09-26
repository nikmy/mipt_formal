package modify

import "mipt_formal/internal/fsm"

type Modifier func(fsm.Machine) (fsm.Machine, error)
