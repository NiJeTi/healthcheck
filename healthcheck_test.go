package healthcheck_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nijeti/healthcheck"
	mocks "github.com/nijeti/healthcheck/internal/generated/mocks" //nolint:revive // path contains different packages
)

func TestNew(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(
		t,
		"healthcheck degradation timeout must be less than unhealthy timeout",
		func() {
			healthcheck.New(
				healthcheck.WithTimeoutDegraded(5*time.Second),
				healthcheck.WithTimeoutUnhealthy(5*time.Second),
			)
		},
	)

	assert.PanicsWithValue(
		t,
		"healthcheck degradation timeout must be less than unhealthy timeout",
		func() {
			healthcheck.New(
				healthcheck.WithTimeoutDegraded(10*time.Second),
				healthcheck.WithTimeoutUnhealthy(5*time.Second),
			)
		},
	)

	assert.NotPanics(
		t, func() {
			healthcheck.New()
		},
	)
}

func TestHealthcheck_Handle(t *testing.T) {
	t.Parallel()

	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	tests := map[string]struct {
		status healthcheck.Status
		ctx    func() context.Context
		setup  func(t *testing.T) *healthcheck.Healthcheck
	}{
		"no_probes": {
			status: healthcheck.StatusUnknown,
			setup: func(_ *testing.T) *healthcheck.Healthcheck {
				return healthcheck.New()
			},
		},
		"context_cancelled": {
			status: healthcheck.StatusUnknown,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				return healthcheck.New(
					healthcheck.WithProbe("probe", mocks.NewMockProbe(t)),
				)
			},
		},
		"one_probe_healthy": {
			status: healthcheck.StatusHealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				probe := mocks.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).Return(nil)

				return healthcheck.New(healthcheck.WithProbe("probe", probe))
			},
		},
		"one_probe_timeout_degraded": {
			status: healthcheck.StatusDegraded,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				probe := mocks.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(20 * time.Millisecond)
						return nil
					},
				)

				return healthcheck.New(
					healthcheck.WithTimeoutDegraded(10*time.Millisecond),
					healthcheck.WithProbe("probe", probe),
				)
			},
		},
		"one_probe_timeout_unhealthy": {
			status: healthcheck.StatusUnhealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				probe := mocks.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(30 * time.Millisecond)
						return nil
					},
				)

				return healthcheck.New(
					healthcheck.WithTimeoutDegraded(10*time.Millisecond),
					healthcheck.WithTimeoutUnhealthy(20*time.Millisecond),
					healthcheck.WithProbe("probe", probe),
				)
			},
		},
		"one_probe_error": {
			status: healthcheck.StatusUnhealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				probe := mocks.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).Return(
					errors.New("probe error"),
				)

				return healthcheck.New(healthcheck.WithProbe("probe", probe))
			},
		},
		"one_probe_panic": {
			status: healthcheck.StatusUnhealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				probe := mocks.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).Panic("probe panic")

				return healthcheck.New(healthcheck.WithProbe("probe", probe))
			},
		},
		"multiple_probes_healthy": {
			status: healthcheck.StatusHealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				p1 := mocks.NewMockProbe(t)
				p2 := mocks.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).Return(nil)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return healthcheck.New(
					healthcheck.WithProbe("p1", p1),
					healthcheck.WithProbe("p2", p2),
				)
			},
		},
		"multiple_probes_one_timeout_degraded": {
			status: healthcheck.StatusDegraded,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				p1 := mocks.NewMockProbe(t)
				p2 := mocks.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(20 * time.Millisecond)
						return nil
					},
				)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return healthcheck.New(
					healthcheck.WithTimeoutDegraded(10*time.Millisecond),
					healthcheck.WithProbe("p1", p1),
					healthcheck.WithProbe("p2", p2),
				)
			},
		},
		"multiple_probes_one_timeout_unhealthy": {
			status: healthcheck.StatusUnhealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				p1 := mocks.NewMockProbe(t)
				p2 := mocks.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(30 * time.Millisecond)
						return nil
					},
				)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return healthcheck.New(
					healthcheck.WithTimeoutDegraded(10*time.Millisecond),
					healthcheck.WithTimeoutUnhealthy(20*time.Millisecond),
					healthcheck.WithProbe("p1", p1),
					healthcheck.WithProbe("p2", p2),
				)
			},
		},
		"multiple_probes_one_error": {
			status: healthcheck.StatusUnhealthy,
			setup: func(t *testing.T) *healthcheck.Healthcheck {
				p1 := mocks.NewMockProbe(t)
				p2 := mocks.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).Return(
					errors.New("p1 error"),
				)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return healthcheck.New(
					healthcheck.WithProbe("p1", p1),
					healthcheck.WithProbe("p2", p2),
				)
			},
		},
	}

	for name, tt := range tests {
		t.Run(
			name, func(t *testing.T) {
				hc := tt.setup(t)

				ctx := context.Background()
				if tt.ctx != nil {
					ctx = tt.ctx()
				}

				status := hc.Handle(ctx)
				assert.Equal(t, tt.status, status)
			},
		)
	}
}
