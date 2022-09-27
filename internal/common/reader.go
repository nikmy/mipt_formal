package common

import (
    "errors"
    "fmt"
    "os"
)

func NewStdinReader() Reader {
    return &reader{in: os.Stdin}
}

type reader struct {
    in *os.File
}

func (r *reader) Read() (string, error) {
    var input string
    _, err := fmt.Fscanf(os.Stdin, "%s", &input)
    if err != nil {
        return "", errors.New(fmt.Sprintf("couldn't read input: %s", err))
    }
    return input, nil
}
