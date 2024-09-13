package healthcheck

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nijeti/healthcheck/internal/generated/mocks"
)

func TestWithLogger(t *testing.T) {
	t.Parallel()

	hc := &Healthcheck{}

	assert.PanicsWithValue(
		t, "healthcheck logger cannot be nil",
		func() {
			WithLogger(nil)(hc)
		},
	)

	logger := slog.Default()
	WithLogger(logger)(hc)
	assert.Equal(t, logger, hc.logger)
}

func TestWithProbe(t *testing.T) {
	t.Parallel()

	probe := healthcheck.NewMockProbe(t)

	hc := &Healthcheck{
		probes: map[string]Probe{},
	}

	assert.PanicsWithValue(
		t, "healthcheck probe cannot be nil", func() {
			WithProbe("probe", nil)(hc)
		},
	)

	WithProbe("probe", probe)(hc)
	assert.Equal(t, probe, hc.probes["probe"])

	assert.PanicsWithValue(
		t, "healthcheck probe 'probe' already registered",
		func() {
			WithProbe("probe", probe)(hc)
		},
	)
}

func TestWithSimpleProbe(t *testing.T) {
	t.Parallel()

	probe := func(_ context.Context) error {
		return nil
	}

	hc := &Healthcheck{
		probes: map[string]Probe{},
	}

	assert.PanicsWithValue(
		t, "healthcheck probe cannot be nil", func() {
			WithSimpleProbe("probe", nil)(hc)
		},
	)

	WithSimpleProbe("probe", probe)(hc)

	assert.PanicsWithValue(
		t, "healthcheck probe 'probe' already registered",
		func() {
			WithSimpleProbe("probe", probe)(hc)
		},
	)
}

func TestWithTimeoutDegraded(t *testing.T) {
	t.Parallel()

	hc := &Healthcheck{}

	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			WithTimeoutDegraded(-1)(hc)
		},
	)

	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			WithTimeoutDegraded(0)(hc)
		},
	)

	timeout := time.Duration(1)
	WithTimeoutDegraded(timeout)(hc)
	assert.Equal(t, timeout, hc.timeoutDegraded)
}

func TestWithTimeoutUnhealthy(t *testing.T) {
	t.Parallel()

	hc := &Healthcheck{}

	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			WithTimeoutUnhealthy(-1)(hc)
		},
	)

	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			WithTimeoutUnhealthy(0)(hc)
		},
	)

	timeout := time.Duration(1)
	WithTimeoutUnhealthy(timeout)(hc)
	assert.Equal(t, timeout, hc.timeoutUnhealthy)
}
