package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	"github.com/server-selfish/backend/internal/presentation"
	"golang.org/x/sync/singleflight"
)

type (
	AuthHandler interface {
		GithubLogin(w http.ResponseWriter, r *http.Request)
		GithubCallback(w http.ResponseWriter, r *http.Request)
		Refresh(w http.ResponseWriter, r *http.Request)
		Me(w http.ResponseWriter, r *http.Request)
		Logout(w http.ResponseWriter, r *http.Request)
	}

	authHandler struct {
		as     service.AuthService
		logger zerolog.Logger
	}
)

var (
	githubCallbackTmpl = template.Must(
		template.New("github_callback.html").ParseFS(
			presentation.PresentationEmbed,
			"templates/github_callback.html",
		),
	)
	refreshGroup singleflight.Group
)

func NewAuthHandler(as service.AuthService, logger zerolog.Logger) AuthHandler {
	return &authHandler{
		as:     as,
		logger: logger,
	}
}

func (h *authHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	loginURL, err := h.as.GetGithubLoginURL(ctx)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "failed to build github login url",
			Error:   err.Error(),
		})
		return
	}

	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

func (h *authHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "invalid callback params",
			Error:   "missing code ",
		})
		return
	}

	tokenPair, err := h.as.HandleGithubCallback(ctx, code, r.UserAgent(), pkg.ReadIP(r))
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "github authentication failed",
			Error:   err.Error(),
		})
		return
	}

	if tokenPair.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_refresh_token",
			Value:    tokenPair.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.RefreshTokenExpiresIn),
		})
	}
	if tokenPair.AccessToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_access_token",
			Value:    tokenPair.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.AccessTokenExpiresIn),
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := githubCallbackTmpl.Execute(w, struct{}{}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.RefreshTokenRequest
	// get refresh token from body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// fallback to cookie if body isn't provided
		req.RefreshToken = ""
	}
	// get refresh token from cookie
	if strings.TrimSpace(req.RefreshToken) == "" {
		if c, err := r.Cookie("selfish_refresh_token"); err == nil {
			req.RefreshToken = c.Value
		}
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		pkg.WriteJSON(w, http.StatusBadRequest, schema.AuthErrorResponse{
			Message: "refresh token is required",
		})
		return
	}

	// get the new access token inside token pair
	v, err, _ := refreshGroup.Do(req.RefreshToken, func() (any, error) {
		return h.as.RefreshAccessToken(ctx, req.RefreshToken, r.UserAgent(), pkg.ReadIP(r))
	})

	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "failed to refresh token",
			Error:   err.Error(),
		})
		return
	}

	tokenPair, ok := v.(schema.AuthTokenPair)
	if !ok {
		// handle the error, e.g.:
		pkg.WriteJSON(w, http.StatusInternalServerError, schema.AuthErrorResponse{
			Message: "internal error",
			Error:   "failed to cast token pair",
		})
		return
	}

	// rotate cookie refresh token
	if tokenPair.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_refresh_token",
			Value:    tokenPair.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.RefreshTokenExpiresIn),
		})
	}

	if tokenPair.AccessToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "selfish_access_token",
			Value:    tokenPair.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(tokenPair.AccessTokenExpiresIn),
		})
	}

	// send back response with access token
	pkg.WriteJSON(w, http.StatusOK, schema.RefreshTokenResponse{
		Message: "token refreshed",
		Data: schema.AuthTokenPair{
			AccessToken:           tokenPair.AccessToken,
			AccessTokenExpiresIn:  tokenPair.AccessTokenExpiresIn,
			RefreshToken:          tokenPair.RefreshToken,
			RefreshTokenExpiresIn: tokenPair.RefreshTokenExpiresIn,
			TokenType:             tokenPair.TokenType,
		},
	})
}

func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user id from context (ingested from middleware)
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	// get user profile
	me, err := h.as.GetMe(ctx, userID)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "failed to resolve user",
			Error:   err.Error(),
		})
		return
	}

	pkg.WriteJSON(w, http.StatusOK, schema.MeResponse{
		Message: "ok",
		Data: schema.UserMeData{
			ID:             me.ID,
			Provider:       me.Provider,
			ProviderUserID: me.ProviderUserID,
			Username:       me.Username,
			Email:          me.Email,
			AvatarURL:      me.AvatarURL,
		},
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.RefreshToken = ""
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		if c, err := r.Cookie("selfish_refresh_token"); err == nil {
			req.RefreshToken = c.Value
		}
	}
	if strings.TrimSpace(req.RefreshToken) != "" {
		// logout service call
		_ = h.as.Logout(ctx, req.RefreshToken)
	}

	// clear cookie regardless of service result
	http.SetCookie(w, &http.Cookie{
		Name:     "selfish_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "selfish_access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") || r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	pkg.WriteJSON(w, http.StatusOK, schema.LogoutResponse{
		Message: "logout success",
	})
}
