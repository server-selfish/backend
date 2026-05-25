package handler

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

func RequireAuth(as service.AuthService, logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			var accessToken string

			if authHeader == "" {
				cookie, err := r.Cookie("selfish_access_token")
				if err != nil {
					logger.Error().Msg(defined_error.ErrMissingAccessToken.Error())
					pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrMissingAccessToken)
					return
				}
				accessToken = strings.TrimSpace(cookie.Value)
			} else {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
					logger.Error().Msg(defined_error.ErrInvaliAuthorizationFormat.Error())
					pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrInvaliAuthorizationFormat)
					return
				}
				accessToken = strings.TrimSpace(parts[1])
			}

			if accessToken == "" {
				logger.Error().Msg(defined_error.ErrEmptyAccessToken.Error())
				pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrMissingAccessToken)
				return
			}

			claims, err := as.ParseAccessToken(accessToken)
			if err != nil {
				logger.Error().Msg(err.Error())
				pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrInvalidOrExpireAccessToken)
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
