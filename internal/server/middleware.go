package server

import (
	"net/http"
	"time"

	"github.com/failfilin/otus_hw/internal/service"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(l *service.RedisLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l == nil {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rr, r)

		handler := r.Method + " " + r.URL.Path
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		l.Log(handler, ip, rr.status, time.Since(start))
	})
}
