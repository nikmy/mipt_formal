package tools

func NewSet[T comparable]() Set[T] {
    return make(Set[T])
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(key T) bool {
    _, has := s[key]
    return has
}

func (s Set[T]) Insert(key T) {
    s[key] = struct{}{}
}

func (s Set[T]) Delete(key T) {
    delete(s, key)
}
