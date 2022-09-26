package modify

import "mipt_formal/internal/fsm"

var _ Modifier = RemoveEpsilon

func RemoveEpsilon(input fsm.Machine) (fsm.Machine, error) {
    panic("not implemented")
}
