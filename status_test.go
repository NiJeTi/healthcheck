package healthcheck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nijeti/healthcheck"
)

func TestStatus_Int(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t, int(healthcheck.StatusHealthy), healthcheck.StatusHealthy.Int(),
	)
	assert.Equal(
		t, int(healthcheck.StatusDegraded), healthcheck.StatusDegraded.Int(),
	)
	assert.Equal(
		t, int(healthcheck.StatusUnhealthy), healthcheck.StatusUnhealthy.Int(),
	)
}

func TestStatus_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		status healthcheck.Status
		want   string
	}{
		"healthy": {
			status: healthcheck.StatusHealthy,
			want:   "healthy",
		},
		"degraded": {
			status: healthcheck.StatusDegraded,
			want:   "degraded",
		},
		"unhealthy": {
			status: healthcheck.StatusUnhealthy,
			want:   "unhealthy",
		},
		"unknown": {
			status: healthcheck.StatusUnknown,
			want:   "unknown",
		},
		"invalid": {
			status: healthcheck.Status(-2),
			want:   "unknown",
		},
	}

	for name, tt := range tests {
		t.Run(
			name, func(t *testing.T) {
				assert.Equal(t, tt.want, tt.status.String())
			},
		)
	}
}
