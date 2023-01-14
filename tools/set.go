package tools

func NewSet[T comparable](keys ...T) Set[T] {
    s := make(Set[T])
    for _, key := range keys {
        s.Insert(key)
    }
    return s
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Size() int {
    return len(s)
}

func (s Set[T]) Has(key T) bool {
    _, has := s[key]
    return has
}

func (s Set[T]) Insert(key T) bool {
    if s.Has(key) {
        return false
    }
    s[key] = struct{}{}
    return true
}

func (s Set[T]) Delete(key T) {
    delete(s, key)
}

func (s Set[T]) Empty() bool {
    return s.Size() == 0
}

func (s Set[T]) AsSlice() []T {
    slice := make([]T, 0, s.Size())
    for x := range s {
        slice = append(slice, x)
    }
    return slice
}
