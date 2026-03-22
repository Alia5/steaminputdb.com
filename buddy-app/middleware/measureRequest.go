package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Alia5/steaminputdb.com/api"
	"github.com/Alia5/steaminputdb.com/api/ctx"
	"github.com/google/uuid"
)

func MeasureRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sw := &api.StatusWriter{ResponseWriter: w, Status: http.StatusOK}
		start := time.Now()

		sw.Header().Set("X-Request-ID", uuid.NewString())
		next.ServeHTTP(sw, r)

		dur := time.Since(start)

		c := r.Context()
		c = context.WithValue(c, ctx.KeyStatusCode, sw.Status)
		c = context.WithValue(c, ctx.KeyDuration, dur.Milliseconds())
		if sw.Error != nil {
			c = context.WithValue(c, ctx.KeyError, sw.Error)
		}
		*r = *r.WithContext(c)
	})
}
