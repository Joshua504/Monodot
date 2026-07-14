package main

import (
	"log"
	"net"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		requestID := newRequestID()

		w.Header().Set("X-Request-ID", requestID)
		r = withRequestID(r, requestID)

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		log.Printf(
			"[%s] %s[%d]%s %s %s | %s | %dB | %v ",
			requestID,
			statusColor(recorder.statusCode),
			recorder.statusCode,
			reset,
			r.Method,
			r.URL.Path,
			clientIP(r),
			recorder.size,
			time.Since(start),
			//r.UserAgent(),
		)
	})
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forward-For"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
