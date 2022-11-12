package common

import (
    "bytes"
    "fmt"
    "io"
    "os"
)

type StdinLineReader struct {
    buf    [][]byte
    cached bool
}

func (r *StdinLineReader) ReadLine() ([]byte, error) {
    if !r.cached {
        err := r.scan()
        if err != nil {
            return nil, err
        }
        r.cached = true
    }

    if len(r.buf) == 0 {
        return nil, io.EOF
    }

    next := r.buf[0]
    r.buf[0] = nil // avoid memory leak
    r.buf = r.buf[1:]
    return next, nil
}

func (r *StdinLineReader) scan() error {
    var s string

    _, err := fmt.Fscanf(os.Stdin, "%s", &s)
    if err != nil {
        return err
    }

    r.buf = bytes.Split([]byte(s), []byte("\n"))
    return nil
}
