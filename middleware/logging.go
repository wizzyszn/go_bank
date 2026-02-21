package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		statusCode:     http.StatusOK,
		ResponseWriter: w,
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := newResponseWriter(w)

		next(wrapped, r)

		duration := time.Since(start)

		log.Printf(
			"%s %s %d %s %d bytes",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			duration,
			wrapped.written,
		)
	}
}

func DetailedLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		log.Printf("→ %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)

		next(wrapped, r)

		duration := time.Since(start)

		log.Printf(
			"← %s %s [%d] %s (%d bytes)",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			duration,
			wrapped.written,
		)
	}
}

func RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := time.Now().UnixNano()
		w.Header().Set("X-Request-ID", strconv.Itoa(int(requestID)))

		// Could also add to context if needed
		// ctx := context.WithValue(r.Context(), "request_id", requestID)
		// r = r.WithContext(ctx)
		next(w, r)
	}
}
