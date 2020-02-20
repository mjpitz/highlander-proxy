package proxy

import (
	"context"
	"io"
)

// Pipe reads from the provider reader and writes data to the provider writer
// until the underlying connections encounter an error.
func Pipe(parent context.Context, reader io.ReadCloser, writer io.WriteCloser) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	for done := false; !done; {
		data := make([]byte, 256)

		n, err := reader.Read(data)
		if err != nil {
			break
		}

		if _, err := writer.Write(data[:n]); err != nil {
			break
		}

		select {
		case <-ctx.Done():
			done = true
		default:
		}
	}

	_ = reader.Close()
	_ = writer.Close()
}
