package proxy

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
)

// Pipe reads from the provider reader and writes data to the provider writer
// until the underlying connections encounter an error.
func Pipe(parent context.Context, reader io.ReadCloser, writer io.WriteCloser) {
	ctx, cancel := context.WithCancel(parent)
	defer func() {
		cancel()
		_ = reader.Close()
		_ = writer.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := make([]byte, 256)

			n, err := reader.Read(data)
			if err != nil {
				logrus.Tracef("encountered error reading from connection: %v", err)
				return
			}

			if _, err := writer.Write(data[:n]); err != nil {
				logrus.Tracef("encountered error writing to connection: %v", err)
				return
			}
		}
	}
}
