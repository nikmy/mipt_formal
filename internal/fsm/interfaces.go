package fsm

type Machine interface {
    Go(from State, by Word) []State
    DOA() string
}

type Modifier interface {
    Modify(Machine) Machine
}
