package http

import (
	"net"
	"strings"

	"github.com/nijeti/healthcheck"
)

// Option configures a Healthcheck server instance.
type Option func(server *Server)

// WithListener sets a custom net.Listener for the Healthcheck server.
// Panics if listener is empty.
func WithListener(listener net.Listener) Option {
	if listener == nil {
		panic("healthcheck server listener cannot be nil")
	}

	return func(s *Server) {
		s.listen = func() (net.Listener, error) {
			return listener, nil
		}
	}
}

// WithAddress sets the address for the Healthcheck server.
// Panics if address is empty.
func WithAddress(address string) Option {
	if address == "" {
		panic("healthcheck server address cannot be empty")
	}

	return func(server *Server) {
		server.listen = listen(address)
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

// WithStatusAdapter sets a custom adapter function
// for converting healthcheck status.
// Panics if adapterFunc is nil.
func WithStatusAdapter(
	adapterFunc func(status healthcheck.Status) (code int, message string),
) Option {
	if adapterFunc == nil {
		panic("healthcheck server status adapter func cannot be nil")
	}

	return func(server *Server) {
		server.statusAdapterFunc = adapterFunc
	}
}
