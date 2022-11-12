package doa

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func NewStdoutWriter() Writer {
	return &doaWriter{
		out: os.Stdout,
	}
}

func NewWriter(w io.Writer) Writer {
	return &doaWriter{
		out: w,
	}
}

func NewFileWriter(filename string) (Writer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("couldn't open output file: %s", err))
	}
	return &doaWriter{out: file}, nil
}

type doaWriter struct {
	out io.Writer
}

func (w *doaWriter) Write(a serializable) error {
	_, err := w.out.Write([]byte(a.DOA()))
	if err != nil {
		return errors.New(fmt.Sprintf("couldn't write doa: %s", err))
	}
	return nil
}
