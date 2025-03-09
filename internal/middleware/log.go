package middleware

import (
	"net/http"

	"github.com/lcafe/a_simple_webserver/internal/logs"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.LogRequest(r)

		lrw := &logs.LogResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		logs.LogResponse(lrw)
	})
}
