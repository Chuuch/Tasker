package metrics

import (
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	return w.ResponseWriter.Write(body)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		RequestStarted()
		defer RequestFinished()

		recorder := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		path := r.Pattern
		if path == "" {
			path = r.URL.Path
		}

		if len(path) > 0 {
			if i := strings.IndexByte(path, ' '); i >= 0 {
				path = path[i+1:]
			}
		}

		ObserveRequest(
			r.Method,
			r.Pattern,
			recorder.status,
			time.Since(start),
		)
	})
}
