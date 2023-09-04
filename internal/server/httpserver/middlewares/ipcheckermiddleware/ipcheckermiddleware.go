// Package: ipcheckermiddleware provides a middleware for checking the IP address of the request.
// IP address taken from the X-Real-IP header.
package ipcheckermiddleware

import (
	"go.uber.org/zap"
	"net"
	"net/http"
)

func IPCheckMiddleware(logger *zap.Logger, subnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if subnet != nil {
				ip := r.Header.Get("X-Real-IP")
				if ip == "" {
					logger.Error("X-Real-IP header is empty")
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				remoteIP := net.ParseIP(ip)
				if !subnet.Contains(remoteIP) {
					logger.Error("IP is not in trust subnet", zap.String("ip", ip))
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
