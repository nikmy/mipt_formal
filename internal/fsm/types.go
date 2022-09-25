package fsm

type State uint64
type Word string

const Epsilon Word = "@"

type Transition struct {
    From State
    To   State
    By   Word
}
