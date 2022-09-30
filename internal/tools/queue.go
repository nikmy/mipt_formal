package tools

import (
    "github.com/gammazero/deque"
)

func NewQueue[T any](elems ...T) *Queue[T] {
    q := &Queue[T]{
        data: deque.New[T](),
    }
    for _, elem := range elems {
        q.Push(elem)
    }
    return q
}

// Queue is interface-simplified wrapper of Deque
type Queue[T any] struct {
    data *deque.Deque[T]
}

func (q *Queue[T]) Head() T {
    return q.data.Front()
}

func (q *Queue[T]) Tail() T {
    return q.data.Back()
}

func (q *Queue[T]) Empty() bool {
    return q.Size() == 0
}

func (q *Queue[T]) Size() int {
    return q.data.Len()
}

func (q *Queue[T]) Push(x T) {
    q.data.PushFront(x)
}

func (q *Queue[T]) Pop() T {
    return q.data.PopBack()
}
