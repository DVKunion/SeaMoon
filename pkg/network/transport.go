package network

import (
	"io"
)

const bufferSize = 64 * 1024

// Transport rw1 and rw2
func Transport(rw1, rw2 io.ReadWriter) error {
	errC := make(chan error, 1)
	go func() {
		errC <- CopyBuffer(rw1, rw2, bufferSize)
	}()

	go func() {
		errC <- CopyBuffer(rw2, rw1, bufferSize)
	}()

	if err := <-errC; err != nil && err != io.EOF {
		return err
	}

	return nil
}

func CopyBuffer(dst io.Writer, src io.Reader, bufSize int) error {
	buf := GetBuffer(bufSize)
	defer PutBuffer(buf)

	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
