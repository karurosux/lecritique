package middleware

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	window   time.Duration
}

type visitor struct {
	lastSeen time.Time
	count    int
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			
			rl.mu.Lock()
			v, exists := rl.visitors[ip]
			if !exists {
				rl.visitors[ip] = &visitor{
					lastSeen: time.Now(),
					count:    1,
				}
				rl.mu.Unlock()
				return next(c)
			}

			if time.Since(v.lastSeen) > rl.window {
				v.count = 1
				v.lastSeen = time.Now()
				rl.mu.Unlock()
				return next(c)
			}

			if v.count >= rl.rate {
				rl.mu.Unlock()
				return response.Error(c, errors.New("RATE_LIMIT", "Too many requests", 429))
			}

			v.count++
			v.lastSeen = time.Now()
			rl.mu.Unlock()

			return next(c)
		}
	}
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
