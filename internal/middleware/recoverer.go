package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"go-srv-temp/internal/httperr"
)

func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"panic", rec,
						"stack", string(debug.Stack()),
					)
					httperr.RespondError(w, httperr.ErrInternal)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
