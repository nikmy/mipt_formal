package common

import (
    "fmt"
    "log"
)

func NewLogger() Logger {
    return logger{}
}

type logger struct{}

func (logger) Fatal(err error) {
    log.Fatalf("[ERROR] %s", err)
}

func (logger) Info(msg string) {
    log.Printf("[INFO] %s", msg)
}

func (l logger) Infof(format string, args ...any) {
    l.Info(fmt.Sprintf(format, args...))
}
