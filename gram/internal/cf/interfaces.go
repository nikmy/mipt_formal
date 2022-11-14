package cf

type reader interface {
    ReadLine() ([]byte, error)
}
