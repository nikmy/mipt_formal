package cfgram

type reader interface {
    ReadLine() ([]byte, error)
}
