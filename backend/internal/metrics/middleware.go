package metrics

import (
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

func (w *responseWriter) unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldSkipMetrics(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		RequestStarted()
		defer RequestFinished()

		recorder := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		ObserveRequest(
			r.Method,
			normalizeRoute(r.Pattern, r.URL.Path),
			recorder.status,
			time.Since(start),
		)
	})
}
func shouldSkipMetrics(path string) bool {
	return path == "/health" || path == "/metrics"
}

func normalizeRoute(pattern, rawPath string) string {
	if pattern == "" {
		return "unmatched"
	}

	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		return pattern[i+1:]
	}
	return pattern
}
