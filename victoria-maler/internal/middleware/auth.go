package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"ecochain-victoria/internal/utils"
)

const (
	cookieName     = "token"
	headerAuth     = "Authorization"
	headerPrefix   = "Bearer "
)

type contextKey string

const userContextKey contextKey = "claims"

// Claims extracts the authenticated claims from the request context.
func Claims(r *http.Request) (*utils.Claims, bool) {
	claims, ok := r.Context().Value(userContextKey).(*utils.Claims)
	return claims, ok
}

func newContextWithClaims(r *http.Request, claims *utils.Claims) context.Context {
	return context.WithValue(r.Context(), userContextKey, claims)
}

// RequireAuth validates the JWT (cookie or Authorization header) and
// stores the claims on the request context. Unauthenticated requests get 401.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := tokenFromRequest(r)
			if err != nil {
				utils.Unauthorized(w, "authentication required")
				return
			}
			claims, err := utils.ParseToken(secret, token)
			if err != nil {
				utils.Unauthorized(w, "invalid or expired token")
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func tokenFromRequest(r *http.Request) (string, error) {
	if header := r.Header.Get(headerAuth); strings.HasPrefix(header, headerPrefix) {
		return strings.TrimPrefix(header, headerPrefix), nil
	}
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	return "", errors.New("no token found")
}
