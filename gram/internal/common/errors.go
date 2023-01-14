package common

import (
    "errors"
    "fmt"
)

var (
    IsNotNonTerminalSymbolError = errors.New("symbol is not non-terminal")
)

func Wrap(err error, wrapper string) error {
    return errors.New(fmt.Sprintf("%s: %s", wrapper, err))
}

func Wrapf(err error, format string, args ...any) error {
    return Wrap(err, fmt.Sprintf(format, args...))
}

func Error(msg string) error {
    return errors.New(msg)
}

func Errorf(format string, args ...any) error {
    return errors.New(fmt.Sprintf(format, args...))
}
