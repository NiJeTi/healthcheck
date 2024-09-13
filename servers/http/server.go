package http

import (
	"log/slog"
	"net/http"

	"github.com/nijeti/healthcheck"
)

// Server represents an HTTP server
// based on net/http package for running health checks.
type Server struct {
	hc                *healthcheck.Healthcheck
	server            *http.Server
	logger            *slog.Logger
	address           string
	route             string
	statusAdapterFunc func(status healthcheck.Status) (int, string)
}

// New creates a new Server instance
// operating provided Healthcheck instance and with the provided options.
func New(hc *healthcheck.Healthcheck, opts ...Option) *Server {
	s := &Server{
		hc:                hc,
		logger:            slog.Default(),
		address:           ":8080",
		route:             "/health",
		statusAdapterFunc: defaultAdapter,
	}

	for _, opt := range opts {
		opt(s)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(s.route, s.handle)
	s.server = &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	return s
}

// Start launches the HTTP server in a separate goroutine to handle health check requests.
func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
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
