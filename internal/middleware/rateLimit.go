package middleware

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/Adega318/urlittle/internal/config"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	conf    config.RateLimitingConfig
	clients *lru.Cache[string, *rate.Limiter]
}

func NewClientLimiter(conf config.RateLimitingConfig) (*clientLimiter, error) {
	clientCache, err := lru.New[string, *rate.Limiter](conf.UserListSize)
	if err != nil {
		return nil, fmt.Errorf("client limiter initialitation error: %w", err)
	}

	return &clientLimiter{conf, clientCache}, nil
}

func (c *clientLimiter) get(ip string) *rate.Limiter {
	if limiter, ok := c.clients.Get(ip); ok {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(c.conf.Limit), c.conf.Burst)
	c.clients.Add(ip, limiter)
	return limiter
}

func (c *clientLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !c.get(ip).Allow() {
			w.Header().Set("Retry-After", "1")
			slog.Debug("rate limiting", "ip", ip)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
