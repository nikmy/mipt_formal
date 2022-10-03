package tools

func readString(s string) *stringReader {
    return &stringReader{
        buf: s,
    }
}

type stringReader struct {
    buf string
}

func (s stringReader) Read() (string, error) {
    return s.buf, nil
}

type stringWriter struct {
    buf string
}

func (s *stringWriter) Write(data []byte) (int, error) {
    s.buf = string(data)
    return 0, nil
}

func (s *stringWriter) String() string {
    return s.buf
}
