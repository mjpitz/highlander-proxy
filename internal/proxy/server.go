package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mjpitz/highlander-proxy/internal/election"
)

// Server defines a leader-aware proxy server.
type Server struct {
	Protocol      string
	BindAddress   string
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
	if target == s.BindAddress {
		target = s.RemoteAddress
	}

	return target, ctx, nil
}

func (s *Server) handleConnection(client net.Conn) {
	target, parent, err := s.GetForwardAddress()
	if err != nil {
		// no leader
		_ = client.Close()
		return
	}

	log.Println("forwarding to", target)

	server, err := net.Dial(s.Protocol, target)
	if err != nil {
		// backend not reachable
		_ = client.Close()
		return
	}

	go pipe(parent, client, server)
	go pipe(parent, server, client)
}

// Serve creates a listener for the configured protocol (tcp or udp)
func (s *Server) Serve() error {
	listener, err := net.Listen(s.Protocol, s.BindAddress)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		_ = listener.Close()
	}()

	for {
		accepted, err := listener.Accept()
		if err != nil {
			// server shutdown
			return nil
		}

		go s.handleConnection(accepted)
	}
}
