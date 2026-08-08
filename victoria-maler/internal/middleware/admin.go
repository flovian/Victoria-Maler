package middleware

import (
	"net/http"

	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/utils"
)

// RequireAdmin restricts a route to authenticated admin users.
func RequireAdmin(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := Claims(r)
			if !ok {
				utils.Unauthorized(w, "authentication required")
				return
			}
			if claims.Role != models.RoleAdmin {
				utils.Error(w, http.StatusForbidden, "admin access required")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
