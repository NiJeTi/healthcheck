package healthcheck

import (
	"context"
)

// Probe defines an interface for performing health checks.
type Probe interface {
	// Check performs a health check and returns an error if the check fails.
	Check(ctx context.Context) error
}

// SimpleProbe defines a function type for performing a health check.
type SimpleProbe func(ctx context.Context) error

type simpleProbe struct {
	check SimpleProbe
}

func (p *simpleProbe) Check(ctx context.Context) error {
	return p.check(ctx)
}
