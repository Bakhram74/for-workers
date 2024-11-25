package logger

import (
	"net/http"
	"time"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

func (lw *logResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := newLoggingResponseWriter(w)

		next.ServeHTTP(lw, r)

		Log = logCallerDesable()

		Log.Info().
			Str("method", r.Method).
			Str("url", r.URL.RequestURI()).
			Str("user_agent", r.UserAgent()).
			Dur("elapsed_ms", time.Since(start)).
			Int("status_code", lw.statusCode).
			Msg("incoming request")
	})
}
