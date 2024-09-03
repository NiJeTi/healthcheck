package healthcheck

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nijeti/healthcheck/internal/generated/mocks/github.com/nijeti/healthcheck"
)

func TestNew(t *testing.T) {
	assert.PanicsWithValue(
		t,
		"healthcheck degradation timeout must be less than unhealthy timeout",
		func() {
			New(
				WithTimeoutDegraded(5*time.Second),
				WithTimeoutUnhealthy(5*time.Second),
			)
		},
	)

	assert.PanicsWithValue(
		t,
		"healthcheck degradation timeout must be less than unhealthy timeout",
		func() {
			New(
				WithTimeoutDegraded(10*time.Second),
				WithTimeoutUnhealthy(5*time.Second),
			)
		},
	)

	assert.NotPanics(
		t, func() {
			New()
		},
	)
}

func TestHealthcheck_Handle(t *testing.T) {
	tests := map[string]struct {
		status Status
		ctx    func() context.Context
		setup  func(t *testing.T) *Healthcheck
	}{
		"no_probes": {
			status: StatusUnknown,
			setup: func(t *testing.T) *Healthcheck {
				return New()
			},
		},
		"context_cancelled": {
			status: StatusUnknown,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setup: func(t *testing.T) *Healthcheck {
				return New(WithProbe("probe", healthcheck.NewMockProbe(t)))
			},
		},
		"one_probe_healthy": {
			status: StatusHealthy,
			setup: func(t *testing.T) *Healthcheck {
				probe := healthcheck.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).Return(nil)

				return New(WithProbe("probe", probe))
			},
		},
		"one_probe_timeout_degraded": {
			status: StatusDegraded,
			setup: func(t *testing.T) *Healthcheck {
				probe := healthcheck.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(20 * time.Millisecond)
						return nil
					},
				)

				return New(
					WithTimeoutDegraded(10*time.Millisecond),
					WithProbe("probe", probe),
				)
			},
		},
		"one_probe_timeout_unhealthy": {
			status: StatusUnhealthy,
			setup: func(t *testing.T) *Healthcheck {
				probe := healthcheck.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(30 * time.Millisecond)
						return nil
					},
				)

				return New(
					WithTimeoutDegraded(10*time.Millisecond),
					WithTimeoutUnhealthy(20*time.Millisecond),
					WithProbe("probe", probe),
				)
			},
		},
		"one_probe_error": {
			status: StatusUnhealthy,
			setup: func(t *testing.T) *Healthcheck {
				probe := healthcheck.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).Return(
					errors.New("probe error"),
				)

				return New(WithProbe("probe", probe))
			},
		},
		"one_probe_panic": {
			status: StatusUnhealthy,
			setup: func(t *testing.T) *Healthcheck {
				probe := healthcheck.NewMockProbe(t)
				probe.EXPECT().Check(mock.Anything).Panic("probe panic")

				return New(WithProbe("probe", probe))
			},
		},
		"multiple_probes_healthy": {
			status: StatusHealthy,
			setup: func(t *testing.T) *Healthcheck {
				p1 := healthcheck.NewMockProbe(t)
				p2 := healthcheck.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).Return(nil)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return New(
					WithProbe("p1", p1),
					WithProbe("p2", p2),
				)
			},
		},
		"multiple_probes_one_timeout_degraded": {
			status: StatusDegraded,
			setup: func(t *testing.T) *Healthcheck {
				p1 := healthcheck.NewMockProbe(t)
				p2 := healthcheck.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(20 * time.Millisecond)
						return nil
					},
				)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return New(
					WithTimeoutDegraded(10*time.Millisecond),
					WithProbe("p1", p1),
					WithProbe("p2", p2),
				)
			},
		},
		"multiple_probes_one_timeout_unhealthy": {
			status: StatusUnhealthy,
			setup: func(t *testing.T) *Healthcheck {
				p1 := healthcheck.NewMockProbe(t)
				p2 := healthcheck.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).RunAndReturn(
					func(_ context.Context) error {
						time.Sleep(30 * time.Millisecond)
						return nil
					},
				)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return New(
					WithTimeoutDegraded(10*time.Millisecond),
					WithTimeoutUnhealthy(20*time.Millisecond),
					WithProbe("p1", p1),
					WithProbe("p2", p2),
				)
			},
		},
		"multiple_probes_one_error": {
			status: StatusUnhealthy,
			setup: func(t *testing.T) *Healthcheck {
				p1 := healthcheck.NewMockProbe(t)
				p2 := healthcheck.NewMockProbe(t)
				p1.EXPECT().Check(mock.Anything).Return(
					errors.New("p1 error"),
				)
				p2.EXPECT().Check(mock.Anything).Return(nil)

				return New(
					WithProbe("p1", p1),
					WithProbe("p2", p2),
				)
			},
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt

		t.Run(
			name, func(t *testing.T) {
				t.Parallel()
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
