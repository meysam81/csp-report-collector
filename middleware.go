package main

import (
	"fmt"
	"net/http"

	"github.com/meysam81/x/ratelimit"
)

var (
	ErrTooManyRequests = []byte(`Too many requests`)
)

func (a *AppState) RateLimitMiddleware(next http.Handler) http.Handler {
	rl := ratelimit.RateLimit{
		Redis:       a.redisClient,
		MaxRequests: a.config.RateLimitMaxRPS,
		RefillRate:  a.config.RateLimitRefillRate,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if xff := r.Header.Get("x-forwarded-for"); xff != "" {
			clientIP = xff
		} else if realip := r.Header.Get("x-real-ip"); realip != "" {
			clientIP = realip
		}

		rate := rl.TokenBucket(r.Context(), clientIP)

		w.Header().Set("x-ratelimit-total", fmt.Sprintf("%d", rate.Total))
		w.Header().Set("x-ratelimit-remaining", fmt.Sprintf("%d", rate.Remaining))

		if !rate.Allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Header().Set("retry-after", fmt.Sprintf("%d", rate.ResetAt().Second()))
			_, err := w.Write(ErrTooManyRequests)
			if err != nil {
				a.logger.Error().Err(err).Msg("failed writing response body")
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
