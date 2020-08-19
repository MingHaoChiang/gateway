package common

import (
	"log"
	"net/http"
	"strings"
	"time"
)

//Logger is a function for log message
func Logger(inner http.Handler, name string, pattern string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)
		responseTime := time.Since(start)
		if strings.Index(r.RequestURI, "metrics") < 0 {
			log.SetOutput(WebLog)
			log.Printf(
				"%s %s %s %s %s %s",
				r.Header.Get("X-Forwarded-For"),
				r.Method,
				r.RequestURI,
				name,
				r.UserAgent(),
				responseTime,
			)
		}
	})
}
