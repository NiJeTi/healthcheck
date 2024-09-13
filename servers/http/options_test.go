package http

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nijeti/healthcheck"
)

func TestWithLogger(t *testing.T) {
	t.Parallel()

	s := &Server{}

	assert.PanicsWithValue(
		t, "healthcheck server logger cannot be nil",
		func() {
			WithLogger(nil)(s)
		},
	)

	logger := slog.Default()
	WithLogger(logger)(s)
	assert.Equal(t, logger, s.logger)
}

func TestWithAddress(t *testing.T) {
	t.Parallel()

	s := &Server{}

	assert.PanicsWithValue(
		t, "healthcheck server address cannot be empty",
		func() {
			WithAddress("")(s)
		},
	)

	address := "127.0.0.1:8080"
	WithAddress(address)(s)
	assert.Equal(t, address, s.address)
}

func TestWithRoute(t *testing.T) {
	t.Parallel()

	s := &Server{}

	assert.PanicsWithValue(
		t, "healthcheck server route is invalid",
		func() {
			WithRoute("")(s)
		},
	)
	assert.PanicsWithValue(
		t, "healthcheck server route is invalid",
		func() {
			WithRoute("health")(s)
		},
	)

	route := "/health"
	WithRoute(route)(s)
	assert.Equal(t, route, s.route)
}

func TestWithStatusAdapter(t *testing.T) {
	t.Parallel()

	s := &Server{}

	assert.PanicsWithValue(
		t, "healthcheck server status adapter func cannot be nil",
		func() {
			WithStatusAdapter(nil)(s)
		},
	)

	adapter := func(status healthcheck.Status) (int, string) {
		return 200, "OK"
	}
	WithStatusAdapter(adapter)(s)

	wantCode, wantMsg := adapter(0)
	gotCode, gotMsg := s.statusAdapterFunc(0)
	assert.Equal(t, wantCode, gotCode)
	assert.Equal(t, wantMsg, gotMsg)
}
