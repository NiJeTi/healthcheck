package http

import (
	"net/http"

	"github.com/nijeti/healthcheck"
)

func Handle(hc *healthcheck.Healthcheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)

			_, err := w.Write([]byte("method not allowed"))
			if err != nil {
				hc.Logger().Error("failed to write response", "error", err)
			}

			return
		}

		ctx := r.Context()

		status := hc.Handle(ctx)
		code, message := defaultAdapter(status)

		w.WriteHeader(code)

		_, err := w.Write([]byte(message))
		if err != nil {
			hc.Logger().ErrorContext(
				ctx, "failed to write response", "error", err,
			)
		}
	}
}

func defaultAdapter(status healthcheck.Status) (code int, message string) {
	message = status.String()

	code = http.StatusOK
	if status > healthcheck.StatusHealthy {
		code = http.StatusServiceUnavailable
	}

	return code, message
}
