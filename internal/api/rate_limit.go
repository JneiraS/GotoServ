package api

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitVisitor struct {
	windowStart  time.Time
	requestCount int
}

type rateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	visitors map[string]*rateLimitVisitor
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:    limit,
		window:   window,
		visitors: make(map[string]*rateLimitVisitor),
	}
}

func (limiter *rateLimiter) allow(ip string, now time.Time) bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	for key, visitor := range limiter.visitors {
		if now.Sub(visitor.windowStart) > 2*limiter.window {
			delete(limiter.visitors, key)
		}
	}

	visitor, ok := limiter.visitors[ip]
	if !ok || now.Sub(visitor.windowStart) >= limiter.window {
		limiter.visitors[ip] = &rateLimitVisitor{
			windowStart:  now,
			requestCount: 1,
		}
		return true
	}

	if visitor.requestCount >= limiter.limit {
		return false
	}

	visitor.requestCount++
	return true
}

func rateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(limit, window)

	return func(c *gin.Context) {
		if c.FullPath() == "/health" {
			c.Next()
			return
		}

		if !limiter.allow(remoteIP(c.Request), time.Now()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}

func remoteIP(request *http.Request) string {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}

	return host
}
