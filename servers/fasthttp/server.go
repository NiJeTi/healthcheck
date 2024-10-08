package fasthttp

import (
	"log/slog"
	"net"

	"github.com/valyala/fasthttp"

	"github.com/nijeti/healthcheck"
)

const (
	defaultAddr  = ":8080"
	defaultRoute = "/health"
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

	s.server = &fasthttp.Server{
		Handler:                      s.handle,
		ErrorHandler:                 s.handleError,
		GetOnly:                      true,
		DisablePreParseMultipartForm: true,
		NoDefaultServerHeader:        true,
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
	err := s.server.Shutdown()
	if err != nil {
		s.logger.Error("failed to stop healthcheck server", "error", err)
	}
}

func (s *Server) handle(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != "/health" {
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

func defaultAdapter(status healthcheck.Status) (int, string) {
	code := fasthttp.StatusInternalServerError
	message := status.String()

	switch status {
	case healthcheck.StatusHealthy:
		code = fasthttp.StatusOK
	case healthcheck.StatusDegraded:
		code = fasthttp.StatusOK
	case healthcheck.StatusUnhealthy:
		code = fasthttp.StatusServiceUnavailable
	}

	return code, message
}
