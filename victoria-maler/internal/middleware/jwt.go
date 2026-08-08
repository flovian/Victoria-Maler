package middleware

import (
	"net/http"

	"ecochain-victoria/internal/utils"
)

// OptAuth is a soft-auth wrapper: it attaches claims when a valid token is
// present but still lets the request through otherwise. Useful for pages
// that show different content to logged-in visitors.
func OptAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token, err := tokenFromRequest(r); err == nil {
				if claims, err := utils.ParseToken(secret, token); err == nil {
					ctx := newContextWithClaims(r, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
