package loggermiddleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogMiddleware(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fields := []zapcore.Field{
				zap.String("method:", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
			}
			lw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t := time.Now()
			defer func() {
				fields = append(fields, zap.Duration("duration", time.Since(t)))
				fields = append(fields, zap.Int("status", lw.Status()))
				log.Info("request", fields...)

			}()
			next.ServeHTTP(lw, r)
		}
		return http.HandlerFunc(fn)
	}
}
