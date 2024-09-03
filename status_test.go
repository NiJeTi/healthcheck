package healthcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_Int(t *testing.T) {
	assert.Equal(t, int(StatusHealthy), StatusHealthy.Int())
	assert.Equal(t, int(StatusDegraded), StatusDegraded.Int())
	assert.Equal(t, int(StatusUnhealthy), StatusUnhealthy.Int())
}

func TestStatus_String(t *testing.T) {
	tests := map[string]struct {
		status Status
		want   string
	}{
		"healthy": {
			status: StatusHealthy,
			want:   "healthy",
		},
		"degraded": {
			status: StatusDegraded,
			want:   "degraded",
		},
		"unhealthy": {
			status: StatusUnhealthy,
			want:   "unhealthy",
		},
		"unknown": {
			status: StatusUnknown,
			want:   "unknown",
		},
		"invalid": {
			status: Status(-2),
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
