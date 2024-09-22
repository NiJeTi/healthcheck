package http

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/nijeti/healthcheck"
)

const (
	defaultAddr  = ":8080"
	defaultRoute = "/health"
)

// Server represents an HTTP server
// based on net/http package for running health checks.
type Server struct {
	hc                *healthcheck.Healthcheck
	server            *http.Server
	logger            *slog.Logger
	listen            func() (net.Listener, error)
	route             string
	statusAdapterFunc func(status healthcheck.Status) (int, string)
}

// New creates a new Server instance
// operating provided Healthcheck instance and with the provided options.
func New(hc *healthcheck.Healthcheck, opts ...Option) *Server {
	s := &Server{
		hc:                hc,
		logger:            slog.Default(),
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
	}

	return s
}

// Start launches the HTTP server in a separate goroutine to handle health check requests.
func (s *Server) Start() {
	go func() {
		ln, err := s.listen()
		if err != nil {
			s.logger.Error(
				"failed to start healthcheck server listener", "error", err,
			)
			return
		}

		err = s.server.Serve(ln)
		if err != nil {
			s.logger.Error("healthcheck server error", "error", err)
		}
	}()
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() {
	err := s.server.Close()
	if err != nil {
		s.logger.Error("failed to stop healthcheck server", "error", err)
	}
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		_, err := w.Write([]byte("method not allowed"))
		if err != nil {
			s.logger.Error("failed to write response", "error", err)
		}

		return
	}

	ctx := r.Context()

	status := s.hc.Handle(ctx)
	code, message := s.statusAdapterFunc(status)

	w.WriteHeader(code)

	_, err := w.Write([]byte(message))
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to write response", "error", err)
	}
}

func listen(addr string) func() (net.Listener, error) {
	return func() (net.Listener, error) {
		return net.Listen("tcp", addr)
	}
}

func defaultAdapter(status healthcheck.Status) (int, string) {
	code := http.StatusInternalServerError
	message := status.String()

	switch status {
	case healthcheck.StatusHealthy:
		code = http.StatusOK
	case healthcheck.StatusDegraded:
		code = http.StatusOK
	case healthcheck.StatusUnhealthy:
		code = http.StatusServiceUnavailable
	}

	return code, message
}
