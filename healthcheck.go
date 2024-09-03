package healthcheck

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Healthcheck represents an application health checker with configurable probes and timeouts.
type Healthcheck struct {
	logger           *slog.Logger
	probes           map[string]Probe
	timeoutDegraded  time.Duration
	timeoutUnhealthy time.Duration
}

// New creates a new Healthcheck instance with the provided options.
func New(opts ...Option) *Healthcheck {
	hc := &Healthcheck{
		logger:           slog.Default(),
		probes:           map[string]Probe{},
		timeoutDegraded:  1 * time.Second,
		timeoutUnhealthy: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(hc)
	}

	if hc.timeoutDegraded >= hc.timeoutUnhealthy {
		panic("healthcheck degradation timeout must be less than unhealthy timeout")
	}

	return hc
}

func (hc *Healthcheck) Handle(ctx context.Context) Status {
	probeCount := len(hc.probes)
	if probeCount == 0 {
		return StatusUnknown
	}

	if ctx.Err() != nil {
		return StatusUnknown
	}

	wg := &sync.WaitGroup{}
	wg.Add(probeCount)

	statuses := make(chan Status, probeCount)

	for name, probe := range hc.probes {
		pl := hc.logger.With("probe", name)
		go hc.probeCheck(pl, ctx, wg, statuses, probe)
	}

	wg.Wait()
	close(statuses)

	return hc.calculateStatus(statuses)
}

func (hc *Healthcheck) probeCheck(
	logger *slog.Logger,
	ctx context.Context,
	wg *sync.WaitGroup,
	statuses chan Status,
	probe Probe,
) {
	defer wg.Done()

	defer func() {
		if err := recover(); err != nil {
			logger.ErrorContext(
				ctx,
				"probe panicked",
				"panic", err,
			)

			statuses <- StatusUnhealthy
		}
	}()

	timedCtx, cancel := context.WithTimeout(ctx, hc.timeoutUnhealthy)
	defer cancel()

	probeTime := time.Now()
	err := probe.Check(timedCtx)
	probeDuration := time.Since(probeTime)

	if err != nil || probeDuration > hc.timeoutUnhealthy {
		logger.ErrorContext(
			ctx,
			"failed to probe",
			"error", err,
			"duration", probeDuration.String(),
		)

		statuses <- StatusUnhealthy
		return
	}

	if probeDuration > hc.timeoutDegraded {
		logger.WarnContext(
			ctx,
			"probe is degraded",
			"duration", probeDuration.String(),
		)

		statuses <- StatusDegraded
		return
	}

	statuses <- StatusHealthy
}

func (hc *Healthcheck) calculateStatus(statuses chan Status) Status {
	status := StatusHealthy
	for s := range statuses {
		if s <= status {
			continue
		}

		status = s
		if status == StatusUnhealthy {
			break
		}
	}
	return status
}
