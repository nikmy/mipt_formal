package tools

func NewSet[T comparable](keys ...T) Set[T] {
    s := make(Set[T])
    for _, key := range keys {
        s.Insert(key)
    }
    return s
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
