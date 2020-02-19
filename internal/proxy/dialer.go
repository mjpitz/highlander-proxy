package proxy

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mjpitz/highlander-proxy/internal/election"
)

type Dialer struct {
	Leader        *election.Leader
	Protocol      string
	Identity      string
	RemoteAddress string
}

func (f *Dialer) Dial() (net.Conn, context.Context, error) {
	leader, ctx, ok := f.Leader.Get()
	if !ok {
		return nil, nil, fmt.Errorf("no leader elected")
	}

	target := leader
	if target == f.Identity {
		target = f.RemoteAddress
	}
	log.Println("forwarding to", target)

	conn, err := net.Dial(f.Protocol, target)
	if err != nil {
		return nil, nil, err
	}

	return conn, ctx, nil
}
