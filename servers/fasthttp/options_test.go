package fasthttp

import (
	"log/slog"
	"net"
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

func TestWithListener(t *testing.T) {
	s := &Server{}

	assert.PanicsWithValue(
		t, "healthcheck server listener cannot be nil",
		func() {
			WithListener(nil)(s)
		},
	)

	ln, err := net.Listen("tcp", "127.0.0.1:1337")
	assert.NoError(t, err)
	defer ln.Close()

	WithListener(ln)(s)

	gotLn, gotErr := s.listen()
	assert.NoError(t, gotErr)

	assert.Equal(t, ln, gotLn)
}

func TestWithAddress(t *testing.T) {
	s := &Server{}

	assert.PanicsWithValue(
		t, "healthcheck server address cannot be empty",
		func() {
			WithAddress("")(s)
		},
	)

	address := "127.0.0.1:1337"
	WithAddress(address)(s)

	ln, err := s.listen()
	defer ln.Close()
	
	assert.NoError(t, err)
	assert.Equal(t, address, ln.Addr().String())
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
