package election

import (
	"context"
	"sync"
	"time"
)

// Config defines a generic configuration for managing leader election.
type Config struct {
	Context       context.Context
	Identity      string
	LockNamespace string
	LockName      string

	LeaseDuration time.Duration
	RenewDeadline time.Duration
	RetryPeriod   time.Duration
}

// NewLeader constructs a leader to be used along with a proxy.Server.
func NewLeader() *Leader {
	return &Leader{
		mu: &sync.Mutex{},
	}
}

// Leader encapsulates logic for tracking the current leader.
type Leader struct {
	mu *sync.Mutex

	ctx     context.Context
	cancel  context.CancelFunc
	current string
}

// Get returns the current leader and their corresponding context.
func (l *Leader) Get() (string, context.Context, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.current == "" {
		return "", nil, false
	}

	return l.current, l.ctx, true
}

// Update cancels the current context, then sets a new leader and context.
func (l *Leader) Update(leader string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cancel != nil {
		l.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())

	l.ctx = ctx
	l.cancel = cancel
	l.current = leader
}
