package handler

import (
	"net/http"
	"strings"

	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
)

func RequireAuth(as service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
					Message: "unauthorized",
					Error:   "missing authorization header",
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
					Message: "unauthorized",
					Error:   "invalid authorization header format",
				})
				return
			}

			accessToken := strings.TrimSpace(parts[1])
			if accessToken == "" {
				pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
					Message: "unauthorized",
					Error:   "empty access token",
				})
				return
			}

			claims, err := as.ParseAccessToken(accessToken)
			if err != nil {
				pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
					Message: "unauthorized",
					Error:   "invalid or expired access token",
				})
				return
			}

			ctx := r.Context()
			ctx = pkg.WithAuthUserID(ctx, claims.UserID)
			ctx = pkg.WithAuthSessionID(ctx, claims.SessionID)
			ctx = pkg.WithAuthProvider(ctx, claims.Provider)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
