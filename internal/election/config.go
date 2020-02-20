package election

import (
	"context"
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
