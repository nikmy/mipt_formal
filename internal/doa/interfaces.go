package doa

type serializable interface {
    DOA() string
}

type Writer interface {
    Write(serializable) error
}
