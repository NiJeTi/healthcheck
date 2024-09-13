package healthcheck

import (
	"context"
)

// Probe defines an interface for performing health checks.
type Probe interface {
	// Check performs a health check and returns an error if the check fails.
	Check(ctx context.Context) error
}

// ProbeFunc defines a function type for performing a health check.
type ProbeFunc func(ctx context.Context) error

type probe struct {
	check ProbeFunc
}

func (p *probe) Check(ctx context.Context) error {
	return p.check(ctx)
}
