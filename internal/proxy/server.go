package proxy

import (
	"context"
	"fmt"
	"net"

	"github.com/mjpitz/highlander-proxy/internal/election"

	"github.com/sirupsen/logrus"
)

// Server defines a leader-aware proxy server.
type Server struct {
	Protocol      string
	Identity      string
	RemoteAddress string
	Leader        *election.Leader
}

// GetForwardAddress returns the current address to forward to based on the current leadership.
func (s *Server) GetForwardAddress() (string, context.Context, error) {
	leader, ctx, ok := s.Leader.Get()
	if !ok {
		return "", nil, fmt.Errorf("no leader elected")
	}

	target := leader
	if target == s.Identity {
		target = s.RemoteAddress
	}

	return target, ctx, nil
}

func (s *Server) handleConnection(client net.Conn) {
	target, parent, err := s.GetForwardAddress()
	if err != nil {
		logrus.Tracef("encountered error forwarding connection: %v", err)
		_ = client.Close()
		return
	}

	logrus.Debugf("forwarding to %s", target)

	server, err := net.Dial(s.Protocol, target)
	if err != nil {
		logrus.Errorf("encountered error dialing %s: %v", target, err)
		_ = client.Close()
		return
	}

	go Pipe(parent, client, server)
	go Pipe(parent, server, client)
}

// Serve creates a listener for the configured protocol (tcp or udp)
func (s *Server) Serve(listener net.Listener) {
	for {
		accepted, err := listener.Accept()
		if err != nil {
			logrus.Debugf("no longer accepting connections: %v", err)
			return
		}

		go s.handleConnection(accepted)
	}
}
