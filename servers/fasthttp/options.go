package fasthttp

import (
	"log/slog"
	"strings"

	"github.com/nijeti/healthcheck"
)

// Option configures a Healthcheck server instance.
type Option func(server *Server)

// WithLogger sets the logger for the Healthcheck server.
// Panics if logger is nil.
func WithLogger(logger *slog.Logger) Option {
	if logger == nil {
		panic("healthcheck server logger cannot be nil")
	}

	return func(s *Server) {
		s.logger = logger
	}
}

// WithAddress sets the address for the Healthcheck server.
// Panics if address is empty.
func WithAddress(address string) Option {
	if address == "" {
		panic("healthcheck server address cannot be empty")
	}

	return func(server *Server) {
		server.address = address
	}
}

// WithRoute sets the route path for the Healthcheck server.
// Panics if route is of invalid format.
func WithRoute(route string) Option {
	if !strings.HasPrefix(route, "/") {
		panic("healthcheck server route is invalid")
	}

	return func(server *Server) {
		server.route = route
	}
}

// WithStatusAdapter sets a custom adapter function for converting healthcheck status.
// Panics if adapterFunc is nil.
func WithStatusAdapter(
	adapterFunc func(status healthcheck.Status) (int, string),
) Option {
	if adapterFunc == nil {
		panic("healthcheck server status adapter func cannot be nil")
	}

	return func(server *Server) {
		server.statusAdapterFunc = adapterFunc
	}
}
