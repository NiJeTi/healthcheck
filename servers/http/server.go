package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/nijeti/healthcheck"
)

// Server represents an HTTP server
// based on net/http package for running health checks.
type Server struct {
	hc                *healthcheck.Healthcheck
	server            *http.Server
	listen            func() (net.Listener, error)
	route             string
	statusAdapterFunc func(status healthcheck.Status) (int, string)
}

const (
	defaultAddr  = ":8080"
	defaultRoute = "/health"
)

// New creates a new Server instance
// operating provided Healthcheck instance and with the provided options.
func New(hc *healthcheck.Healthcheck, opts ...Option) *Server {
	s := &Server{
		hc:                hc,
		listen:            listen(defaultAddr),
		route:             defaultRoute,
		statusAdapterFunc: defaultAdapter,
	}

	for _, opt := range opts {
		opt(s)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(s.route, s.handle)
	s.server = &http.Server{
		Handler: mux,
		//nolint:revive,mnd // time span
		ReadHeaderTimeout: 100 * time.Millisecond,
	}

	return s
}

// Start launches the server in a separate goroutine.
// Logs an error if the server fails to start or encounters an issue.
func (s *Server) Start() {
	go func() {
		if err := s.Serve(); err != nil {
			s.hc.Logger().Error("healthcheck server error", "error", err)
			return
		}
	}()
}

// Serve launches the server in the same goroutine.
// Returns an error if the server fails to start or encounters an issue.
func (s *Server) Serve() error {
	ln, err := s.listen()
	if err != nil {
		return fmt.Errorf(
			"failed to start healthcheck server listener: %w", err,
		)
	}

	err = s.server.Serve(ln)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("healthcheck server error: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the server.
// Logs an error if the server shutdown process fails.
func (s *Server) Stop() {
	if err := s.StopWithContext(context.Background()); err != nil {
		s.hc.Logger().Error("failed to stop healthcheck server", "error", err)
	}
}

// StopWithTimeout gracefully shuts down the server with a provided timeout.
// Returns an error if the server shutdown process fails.
func (s *Server) StopWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.StopWithContext(ctx)
}

// StopWithContext gracefully shuts down the server using the provided context.
// Returns an error if the server shutdown process fails.
func (s *Server) StopWithContext(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop healthcheck server: %w", err)
	}

	return nil
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		_, err := w.Write([]byte("method not allowed"))
		if err != nil {
			s.hc.Logger().Error("failed to write response", "error", err)
		}

		return
	}

	ctx := r.Context()

	status := s.hc.Handle(ctx)
	code, message := s.statusAdapterFunc(status)

	w.WriteHeader(code)

	_, err := w.Write([]byte(message))
	if err != nil {
		s.hc.Logger().ErrorContext(
			ctx, "failed to write response", "error", err,
		)
	}
}

func listen(addr string) func() (net.Listener, error) {
	return func() (net.Listener, error) {
		return net.Listen("tcp", addr)
	}
}
