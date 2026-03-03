package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wizzyszn/go_bank/utils"
)

type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*tokenBucket
	rate    float64
	burst   float64
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

func NewRateLimiter(rate, burst float64) *RateLimiter {
	if rate <= 0 {
		rate = 50.0
	}
	if burst < 0 {
		burst = rate * 2
	}

	return &RateLimiter{
		rate:    rate,
		burst:   float64(burst),
		clients: make(map[string]*tokenBucket, 1024),
	}
}

func (rl *RateLimiter) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		if clientIP == "" {
			next(w, r)
			return
		}

		now := time.Now()
		rl.mu.RLock()
		bucket, exists := rl.clients[clientIP]
		rl.mu.RUnlock()

		if !exists {
			rl.mu.Lock()
			bucket, exists = rl.clients[clientIP]
			if !exists {
				bucket = &tokenBucket{
					tokens: rl.burst,
					last:   now,
				}
				rl.clients[clientIP] = bucket
			}
			rl.mu.Unlock()
		}
		elapsed := now.Sub(bucket.last).Seconds()
		newTokens := elapsed * rl.rate
		bucket.tokens = min(bucket.tokens+newTokens, rl.burst)
		bucket.last = now
		if bucket.tokens < 1 {

			waitSec := (1 - bucket.tokens) / rl.rate
			if waitSec < 1 {
				waitSec = 1
			}

			w.Header().Set("Retry-After", strconv.Itoa(int(waitSec)))
			utils.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded. Try again later.")
			return
		}

		bucket.tokens -= 1

		approxRemaining := int(bucket.tokens)
		if approxRemaining < 0 {
			approxRemaining = 0
		}
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(rl.burst)))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(approxRemaining))

		resetSec := int((rl.burst - bucket.tokens) / rl.rate)
		if resetSec > 0 {
			w.Header().Set("X-RateLimit-Reset", strconv.Itoa(resetSec))
		}

		next(w, r)
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}

	for _, h := range headers {
		if val := r.Header.Get(h); val != "" {
			parts := strings.Split(val, ",")
			if len(parts) > 0 {
				ip := strings.TrimSpace(parts[0])
				if net.ParseIP(ip) != nil {
					return ip
				}
			}
		}
	}

	remoteAddr := r.RemoteAddr

	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	}
	return remoteAddr
}
