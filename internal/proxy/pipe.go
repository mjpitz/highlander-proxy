package proxy

import (
	"context"
	"net"
)

func pipe(parent context.Context, reader, writer net.Conn) {
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
