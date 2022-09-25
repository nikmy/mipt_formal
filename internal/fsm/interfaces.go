package fsm

type Machine interface {
    Go(from State, by Word) []State
}

type Modifier interface {
    Modify(Machine) Machine
}
