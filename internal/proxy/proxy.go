package proxy

import (
	"bufio"
	"context"
	"net"
)

func proxy(ctx context.Context, client, server net.Conn) {
	go pipe(ctx, client, server)
	go pipe(ctx, server, client)
}

func pipe(ctx context.Context, readerConn, writerConn net.Conn) {
	reader := bufio.NewReader(readerConn)
	writer := bufio.NewWriter(writerConn)

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

	_ = readerConn.Close()
	_ = writerConn.Close()
}
