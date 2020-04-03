package proxy

import (
	"context"
	"fmt"
	"net"

	"github.com/mjpitz/highlander-proxy/internal/config"
	"github.com/mjpitz/highlander-proxy/internal/election"

	"github.com/sirupsen/logrus"
)

// Server defines a leader-aware proxy server.
type Server struct {
	Route    *config.Route
	Identity string
	Leader   *election.Leader
}

// GetForwardAddress returns the current address to forward to based on the current leadership.
func (s *Server) GetForwardAddress() (string, string, context.Context, error) {
	leader, ctx, ok := s.Leader.Get()
	if !ok {
		return "", "", nil, fmt.Errorf("no leader elected")
	}

	from := s.Route.From()

	protocol := from.Scheme
	target := fmt.Sprintf("%s:%s", leader, from.Port())

	if leader == s.Identity {
		to := s.Route.To()

		protocol = to.Scheme

		if protocol == "unix" {
			// use path for unix sockets
			target = to.Path
		} else {
			target = to.Host
		}
	}

	return protocol, target, ctx, nil
}

func (s *Server) handleConnection(client net.Conn) {
	protocol, target, parent, err := s.GetForwardAddress()
	if err != nil {
		logrus.Tracef("encountered error forwarding connection: %v", err)
		_ = client.Close()
		return
	}

	logrus.Debugf("forwarding to %s://%s", protocol, target)
	server, err := net.Dial(protocol, target)
	if err != nil {
		logrus.Errorf("encountered error dialing %s: %v", target, err)
		_ = client.Close()
		return
	}

	go Pipe(parent, client, server)
	go Pipe(parent, server, client)
}

// Serve creates a listener for the configured protocol (tcp or udp)
func (s *Server) Serve(listener net.Listener, stopCh chan struct{}) {
	defer listener.Close()

	for done := false; !done; {
		select {
		case <-stopCh:
			done = true
			break
		default:
			accepted, err := listener.Accept()
			if err != nil {
				logrus.Debugf("no longer accepting connections: %v", err)
				return
			}

			go s.handleConnection(accepted)
		}
	}
}

func (s *Server) Start(stopCh chan struct{}) error {
	from := s.Route.From()

	protocol := from.Scheme
	bindAddress := from.Host

	listener, err := net.Listen(protocol, bindAddress)
	if err != nil {
		return err
	}

	logrus.Infof("listening on %s://%s", protocol, bindAddress)
	go s.Serve(listener, stopCh)
	return nil
}
