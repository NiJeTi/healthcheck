package healthcheck

import (
	"fmt"
	"log/slog"
	"time"
)

// Option configures a Healthcheck instance.
type Option func(hc *Healthcheck)

// WithLogger sets the logger for the Healthcheck instance.
// Panics if logger is nil.
func WithLogger(logger *slog.Logger) Option {
	if logger == nil {
		panic("healthcheck logger cannot be nil")
	}

	return func(hc *Healthcheck) {
		hc.logger = logger
	}
}

// WithProbe registers a new health check probe with the given name.
// Panics if probe is nil or a probe with the same name already exists.
func WithProbe(name string, probe Probe) Option {
	if probe == nil {
		panic("healthcheck probe cannot be nil")
	}

	return func(hc *Healthcheck) {
		if _, ok := hc.probes[name]; ok {
			p := fmt.Sprintf("healthcheck probe '%s' already registered", name)
			panic(p)
		}

		hc.probes[name] = probe
	}
}

// WithSimpleProbe registers a simple health check probe under the specified name.
// Panics if probe is nil or a probe with the same name already exists.
func WithSimpleProbe(name string, probe SimpleProbe) Option {
	if probe == nil {
		panic("healthcheck probe cannot be nil")
	}

	return WithProbe(name, &simpleProbe{check: probe})
}

// WithTimeoutDegraded sets the time after which a probe is considered degraded.
// Panics if timeout is less than or equal to 0.
func WithTimeoutDegraded(timeout time.Duration) Option {
	if timeout <= 0 {
		panic("healthcheck timeout must be greater than zero")
	}

	return func(hc *Healthcheck) {
		hc.timeoutDegraded = timeout
	}
}

// WithTimeoutUnhealthy sets the time after which a probe is considered unhealthy.
// Panics if timeout is less than or equal to 0.
func WithTimeoutUnhealthy(timeout time.Duration) Option {
	if timeout <= 0 {
		panic("healthcheck timeout must be greater than zero")
	}

	return func(hc *Healthcheck) {
		hc.timeoutUnhealthy = timeout
	}
}
