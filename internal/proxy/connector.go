package proxy

import (
	"context"
	"net"
)

func Connect(root context.Context, listener net.Listener, connections chan net.Conn) {
	for {
		accepted, err := listener.Accept()
		if err != nil {
			// no longer accepting connections
			return
		}

		connections <- accepted

		select {
		case <-root.Done():
			// no longer running
			return
		default:
		}
	}
}

func Forward(root context.Context, dialer *Dialer, connections chan net.Conn) {
	for {
		select {
		case <-root.Done():
			// no longer running
			return
		case client := <-connections:
			server, ctx, err := dialer.Dial()
			if err != nil {
				// can't dial server
				// no leader
			} else {
				proxy(ctx, client, server)
			}
		}
	}
}
