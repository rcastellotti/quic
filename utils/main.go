package utils

import (
	"bufio"
	"io"
)

type bufferedWriteCloser struct {
	*bufio.Writer
	io.Closer
}

func NewBufferedWriteCloser(writer *bufio.Writer, closer io.Closer) io.WriteCloser {
	return &bufferedWriteCloser{
		Writer: writer,
		Closer: closer,
	}
}
