package healthcheck_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nijeti/healthcheck"
	mocks "github.com/nijeti/healthcheck/internal/generated/mocks" //nolint:revive // path contains different packages
)

func TestWithLogger(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(
		t, "healthcheck logger cannot be nil",
		func() {
			healthcheck.New(
				healthcheck.WithLogger(nil),
			)
		},
	)
	assert.NotPanics(
		t, func() {
			healthcheck.New(
				healthcheck.WithLogger(slog.Default()),
			)
		},
	)
}

func TestWithProbe(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(
		t, "healthcheck probe cannot be nil", func() {
			healthcheck.New(
				healthcheck.WithProbe("probe", nil),
			)
		},
	)
	assert.NotPanics(
		t, func() {
			healthcheck.New(
				healthcheck.WithProbe("probe", mocks.NewMockProbe(t)),
			)
		},
	)
	assert.PanicsWithValue(
		t, "healthcheck probe 'probe' already registered",
		func() {
			healthcheck.New(
				healthcheck.WithProbe("probe", mocks.NewMockProbe(t)),
				healthcheck.WithProbe("probe", mocks.NewMockProbe(t)),
			)
		},
	)
}

func TestWithSimpleProbe(t *testing.T) {
	t.Parallel()

	probeFunc := func(_ context.Context) error {
		return nil
	}

	assert.PanicsWithValue(
		t, "healthcheck probe cannot be nil", func() {
			healthcheck.New(
				healthcheck.WithSimpleProbe("probe", nil),
			)
		},
	)
	assert.NotPanics(
		t, func() {
			healthcheck.New(
				healthcheck.WithSimpleProbe("probe", probeFunc),
			)
		},
	)
	assert.PanicsWithValue(
		t, "healthcheck probe 'probe' already registered",
		func() {
			healthcheck.New(
				healthcheck.WithSimpleProbe("probe", probeFunc),
				healthcheck.WithSimpleProbe("probe", probeFunc),
			)
		},
	)
}

func TestWithTimeoutDegraded(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			healthcheck.New(
				healthcheck.WithTimeoutDegraded(-1),
			)
		},
	)
	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			healthcheck.New(
				healthcheck.WithTimeoutDegraded(0),
			)
		},
	)
	assert.NotPanics(
		t, func() {
			healthcheck.New(
				healthcheck.WithTimeoutDegraded(1),
			)
		},
	)
}

func TestWithTimeoutUnhealthy(t *testing.T) {
	t.Parallel()

	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			healthcheck.New(
				healthcheck.WithTimeoutUnhealthy(-1),
			)
		},
	)
	assert.PanicsWithValue(
		t, "healthcheck timeout must be greater than zero",
		func() {
			healthcheck.New(
				healthcheck.WithTimeoutUnhealthy(0),
			)
		},
	)
	assert.NotPanics(
		t, func() {
			healthcheck.New(
				healthcheck.WithTimeoutDegraded(1),
				healthcheck.WithTimeoutUnhealthy(2),
			)
		},
	)
}
