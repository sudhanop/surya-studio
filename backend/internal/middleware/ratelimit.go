package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	loginLimiters = make(map[string]*ipLimiter)
	loginMu       sync.Mutex
)

func LoginRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		loginMu.Lock()
		lim, ok := loginLimiters[ip]
		if !ok {
			lim = &ipLimiter{limiter: rate.NewLimiter(rate.Every(time.Minute), 5), lastSeen: time.Now()}
			loginLimiters[ip] = lim
		}
		lim.lastSeen = time.Now()
		allowed := lim.limiter.Allow()
		loginMu.Unlock()
		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many login attempts"})
		}
		return c.Next()
	}
}
