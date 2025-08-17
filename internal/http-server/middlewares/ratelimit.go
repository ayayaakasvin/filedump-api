package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"up-down-server/internal/http-server/ctx"
	"up-down-server/internal/models"
)

var (
	expTimeForRateLimit time.Duration = time.Second * 4
)
const (
	ratelimitformatstring = "ratelimit:%d/%s" // where %d is ip -> plan is like SET: ratelimitformatstring -> true
)

func (m *Middlewares) RateLimitMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ctx.CtxUserIDKey).(int)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			models.SendErrorJson(w, http.StatusInternalServerError, "invalid remote address")
			return
		}

		key := fmt.Sprintf(ratelimitformatstring, userId, ip)
		set := m.cache.SetNX(r.Context(), key, "1", expTimeForRateLimit);
		if err := set.Err(); err != nil {
			models.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return 
		}

		if !set.Val() {
			models.SendErrorJson(w, http.StatusTooManyRequests, "rate limit")
			return
		}

		h.ServeHTTP(w, r)
	}
}