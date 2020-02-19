package election

import (
	"context"
	"sync"
	"time"
)

type Config struct {
	Context       context.Context
	Identity      string
	LockNamespace string
	LockName      string

	LeaseDuration time.Duration
	RenewDeadline time.Duration
	RetryPeriod   time.Duration
}

func NewLeader() *Leader {
	return &Leader{
		mu: &sync.Mutex{},
	}
}

type Leader struct {
	mu *sync.Mutex

	ctx     context.Context
	cancel  context.CancelFunc
	current string
}

func (l *Leader) Get() (string, context.Context, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.current == "" {
		return "", nil, false
	}

	return l.current, l.ctx, true
}

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
