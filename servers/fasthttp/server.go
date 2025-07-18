//nolint:gci // another workspace module
package fasthttp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/nijeti/healthcheck"
)

// Server represents an HTTP server
// based on fasthttp package for running health checks.
type Server struct {
	hc                *healthcheck.Healthcheck
	server            *fasthttp.Server
	logger            *slog.Logger
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
		logger:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		listen:            listen(defaultAddr),
		route:             defaultRoute,
		statusAdapterFunc: defaultAdapter,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.server = &fasthttp.Server{
		Handler:                      s.handle,
		ErrorHandler:                 s.handleError,
		GetOnly:                      true,
		DisablePreParseMultipartForm: true,
		NoDefaultServerHeader:        true,
	}

	return s
}

// Start launches the server in a separate goroutine.
// Logs an error if the server fails to start or encounters an issue.
func (s *Server) Start() {
	go func() {
		if err := s.Serve(); err != nil {
			s.logger.Error("healthcheck server error", "error", err)
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
	if err != nil {
		return fmt.Errorf("healthcheck server error: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the server.
// Logs an error if the server shutdown process fails.
func (s *Server) Stop() {
	err := s.server.ShutdownWithContext(context.Background())
	if err != nil {
		s.logger.Error("failed to stop healthcheck server", "error", err)
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
	if err := s.server.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("failed to stop healthcheck server: %w", err)
	}

	return nil
}

func (s *Server) handle(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != s.route {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	if !ctx.IsGet() {
		ctx.Error("method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	status := s.hc.Handle(ctx)
	code, message := s.statusAdapterFunc(status)

	ctx.SetStatusCode(code)
	ctx.SetBodyString(message)
}

func (s *Server) handleError(ctx *fasthttp.RequestCtx, err error) {
	s.logger.ErrorContext(ctx, "healthcheck server error", "error", err)
}

func listen(addr string) func() (net.Listener, error) {
	return func() (net.Listener, error) {
		return net.Listen("tcp", addr)
	}
}

func defaultAdapter(status healthcheck.Status) (code int, message string) {
	message = status.String()

	code = fasthttp.StatusOK
	if status > healthcheck.StatusHealthy {
		code = fasthttp.StatusServiceUnavailable
	}

	return
}
