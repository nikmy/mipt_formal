package common

type Logger interface {
	Fatal(err error)
	Info(msg string)
	Infof(format string, args ...any)
}

type Reader interface {
	Read() (string, error)
}

type Writer interface {
	Write(any) error
}
