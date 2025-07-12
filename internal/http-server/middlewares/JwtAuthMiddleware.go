package middlewares

import (
	"net/http"
	"strings"
	"up-down-server/internal/http-server/ctx"
	"up-down-server/internal/lib/jwttool"
	"up-down-server/internal/models"

	"github.com/redis/go-redis/v9"
)

const (
	AuthorizationHeader = "Authorization"
)

// JWTAuthMiddleware is a middleware for http.HandlerFunc
func (m *Middlewares) JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			unauthorized(w, "authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			unauthorized(w, "authorization header missing")
			return
		}

		claims, err := jwttool.ValidateJWT(tokenString)
		if err != nil {
			unauthorized(w, "failed to validate jwt")
			return
		}

		sessionIdAny, exists := claims["session_id"]
		if !exists {
			unauthorized(w, "session_id missing")
			return
		}

		sessionId := sessionIdAny.(string)

		if _, err := m.cache.Get(r.Context(), sessionId); err == redis.Nil {
			unauthorized(w, "session is expired")
			return
		}

		userIdAny, ok := claims["user_id"]
		if !ok || userIdAny == nil {
			unauthorized(w, "user_id missing")
			return
		}

		userIdInt, err := jwttool.FetchUserID(userIdAny)
		if err != nil {
			unauthorized(w, "user_id is invalid")
			return
		}

		r = ctx.WrapValueIntoRequest(r, ctx.CtxUserIDKey, userIdInt)
		r = ctx.WrapValueIntoRequest(r, ctx.CtxSessionIDKey, sessionId)

		next(w, r)
	}
}

func unauthorized(w http.ResponseWriter, msg string)  {
	models.SendErrorJson(w, http.StatusUnauthorized, "%s", msg)
}