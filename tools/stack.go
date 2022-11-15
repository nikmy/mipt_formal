package tools

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0),
	}
}

func StackFromSlice[T any](slice []T) *Stack[T] {
	s := &Stack[T]{}
	s.data = make([]T, len(slice))
	copy(s.data, slice)
	return s
}

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Size() int {
	return len(s.data)
}

func (s *Stack[T]) Top() T {
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Push(x T) {
	s.data = append(s.data, x)
}

func (s *Stack[T]) Pop() T {
	last := len(s.data) - 1
	pop := s.data[last]
	s.data = s.data[:last]
	return pop
}
